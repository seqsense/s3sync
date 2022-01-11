// Copyright 2019 SEQSENSE, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package s3sync

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
)

// Manager manages the sync operation.
type Manager struct {
	s3             s3iface.S3API
	nJobs          int
	del            bool
	dryrun         bool
	acl            *string
	guessMime      bool
	contentType    *string
	downloaderOpts []func(*s3manager.Downloader)
	uploaderOpts   []func(*s3manager.Uploader)
}

type operation int

const (
	opUpdate operation = iota
	opDelete
)

type fileInfo struct {
	name           string
	err            error
	path           string
	size           int64
	lastModified   time.Time
	singleFile     bool
	existsInSource bool
}

type fileOp struct {
	*fileInfo
	op operation
}

// New returns a new Manager.
func New(sess *session.Session, options ...Option) *Manager {
	m := &Manager{
		s3:        s3.New(sess),
		nJobs:     DefaultParallel,
		guessMime: true,
	}
	for _, o := range options {
		o(m)
	}
	return m
}

// Sync syncs the files between s3 and local disks.
func (m *Manager) Sync(source, dest string) error {
	sourceURL, err := url.Parse(source)
	if err != nil {
		return err
	}

	destURL, err := url.Parse(dest)
	if err != nil {
		return err
	}

	chJob := make(chan func())
	var wg sync.WaitGroup
	for i := 0; i < m.nJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range chJob {
				job()
			}
		}()
	}
	defer func() {
		close(chJob)
		wg.Wait()
	}()

	if isS3URL(sourceURL) {
		sourceS3Path, err := urlToS3Path(sourceURL)
		if err != nil {
			return err
		}
		if isS3URL(destURL) {
			destS3Path, err := urlToS3Path(destURL)
			if err != nil {
				return err
			}
			return m.syncS3ToS3(chJob, sourceS3Path, destS3Path)
		}
		return m.syncS3ToLocal(chJob, sourceS3Path, dest)
	}

	if isS3URL(destURL) {
		destS3Path, err := urlToS3Path(destURL)
		if err != nil {
			return err
		}
		return m.syncLocalToS3(chJob, source, destS3Path)
	}

	return errors.New("local to local sync is not supported")
}

func isS3URL(url *url.URL) bool {
	return url.Scheme == "s3"
}

func (m *Manager) syncS3ToS3(chJob chan func(), sourcePath, destPath *s3Path) error {
	return errors.New("S3 to S3 sync feature is not implemented")
}

func (m *Manager) syncLocalToS3(chJob chan func(), sourcePath string, destPath *s3Path) error {
	wg := &sync.WaitGroup{}
	errs := &multiErr{}
	for source := range filterFilesForSync(
		listLocalFiles(sourcePath), m.listS3Files(destPath), m.del,
	) {
		wg.Add(1)
		source := source
		chJob <- func() {
			defer wg.Done()
			if source.err != nil {
				errs.Append(source.err)
				return
			}
			switch source.op {
			case opUpdate:
				if err := m.upload(source.fileInfo, sourcePath, destPath); err != nil {
					errs.Append(err)
				}
			case opDelete:
				if err := m.deleteRemote(source.fileInfo, destPath); err != nil {
					errs.Append(err)
				}
			}
		}
	}
	wg.Wait()

	return errs.ErrOrNil()
}

// syncS3ToLocal syncs the given s3 path to the given local path.
func (m *Manager) syncS3ToLocal(chJob chan func(), sourcePath *s3Path, destPath string) error {
	wg := &sync.WaitGroup{}
	errs := &multiErr{}
	for source := range filterFilesForSync(
		m.listS3Files(sourcePath), listLocalFiles(destPath), m.del,
	) {
		wg.Add(1)
		source := source
		chJob <- func() {
			defer wg.Done()
			if source.err != nil {
				errs.Append(source.err)
				return
			}
			switch source.op {
			case opUpdate:
				if err := m.download(source.fileInfo, sourcePath, destPath); err != nil {
					errs.Append(err)
				}
			case opDelete:
				if err := m.deleteLocal(source.fileInfo, destPath); err != nil {
					errs.Append(err)
				}
			}
		}
	}
	wg.Wait()

	return errs.ErrOrNil()
}

func (m *Manager) download(file *fileInfo, sourcePath *s3Path, destPath string) error {
	var targetFilename string
	if !strings.HasSuffix(destPath, "/") && file.singleFile {
		// Destination path is not a directory and source is a single file.
		targetFilename = destPath
	} else {
		targetFilename = filepath.Join(destPath, file.name)
	}
	targetDir := filepath.Dir(targetFilename)

	println("Downloading", file.name, "to", targetFilename)
	if m.dryrun {
		return nil
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	writer, err := os.Create(targetFilename)
	if err != nil {
		return err
	}

	defer writer.Close()

	var sourceFile string
	if file.singleFile {
		sourceFile = file.name
	} else {
		// Using filepath.ToSlash for change backslash to slash on Windows
		sourceFile = filepath.ToSlash(filepath.Join(sourcePath.bucketPrefix, file.name))
	}

	_, err = s3manager.NewDownloaderWithClient(
		m.s3,
		m.downloaderOpts...,
	).Download(writer, &s3.GetObjectInput{
		Bucket: aws.String(sourcePath.bucket),
		Key:    aws.String(sourceFile),
	})
	if err != nil {
		return err
	}

	err = os.Chtimes(targetFilename, file.lastModified, file.lastModified)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) deleteLocal(file *fileInfo, destPath string) error {
	var targetFilename string
	if !strings.HasSuffix(destPath, "/") && file.singleFile {
		// Destination path is not a directory and source is a single file.
		targetFilename = destPath
	} else {
		targetFilename = filepath.Join(destPath, file.name)
	}

	println("Deleting", targetFilename)
	if m.dryrun {
		return nil
	}

	return os.Remove(targetFilename)
}

func (m *Manager) upload(file *fileInfo, sourcePath string, destPath *s3Path) error {
	var sourceFilename string
	if file.singleFile {
		sourceFilename = sourcePath
	} else {
		sourceFilename = filepath.Join(sourcePath, file.name)
	}

	destFile := *destPath
	if strings.HasSuffix(destPath.bucketPrefix, "/") || destPath.bucketPrefix == "" || !file.singleFile {
		// If source is a single file and destination is not a directory, use destination URL as is.
		// Using filepath.ToSlash for change backslash to slash on Windows
		destFile.bucketPrefix = filepath.ToSlash(filepath.Join(destPath.bucketPrefix, file.name))
	}

	println("Uploading", file.name, "to", destFile.String())
	if m.dryrun {
		return nil
	}

	var contentType *string
	switch {
	case m.contentType != nil:
		contentType = m.contentType
	case m.guessMime:
		mime, err := mimetype.DetectFile(sourceFilename)
		if err != nil {
			return err
		}
		s := mime.String()
		contentType = &s
	}

	reader, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}

	defer reader.Close()

	_, err = s3manager.NewUploaderWithClient(
		m.s3,
		m.uploaderOpts...,
	).Upload(&s3manager.UploadInput{
		Bucket:      aws.String(destFile.bucket),
		Key:         aws.String(destFile.bucketPrefix),
		ACL:         m.acl,
		Body:        reader,
		ContentType: contentType,
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) deleteRemote(file *fileInfo, destPath *s3Path) error {
	destFile := *destPath
	if strings.HasSuffix(destPath.bucketPrefix, "/") || destPath.bucketPrefix == "" || !file.singleFile {
		// If source is a single file and destination is not a directory, use destination URL as is.
		// Using filepath.ToSlash for change backslash to slash on Windows
		destFile.bucketPrefix = filepath.ToSlash(filepath.Join(destPath.bucketPrefix, file.name))
	}

	println("Deleting", destFile.String())
	if m.dryrun {
		return nil
	}

	_, err := m.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(destFile.bucket),
		Key:    aws.String(destFile.bucketPrefix),
	})
	return err
}

// listS3Files return a channel which receives the file infos under the given s3Path.
func (m *Manager) listS3Files(path *s3Path) chan *fileInfo {
	c := make(chan *fileInfo, 50000) // TODO: revisit this buffer size later

	go func() {
		defer close(c)
		var token *string
		for {
			if token = m.listS3FileWithToken(c, path, token); token == nil {
				break
			}
		}
	}()

	return c
}

// listS3FileWithToken lists (send to the result channel) the s3 files from the given continuation token.
func (m *Manager) listS3FileWithToken(c chan *fileInfo, path *s3Path, token *string) *string {
	list, err := m.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:            &path.bucket,
		Prefix:            &path.bucketPrefix,
		ContinuationToken: token,
	})
	if err != nil {
		sendErrorInfoToChannel(c, err)
		return nil
	}

	for _, object := range list.Contents {
		if strings.HasSuffix(*object.Key, "/") {
			// Skip directory like object
			continue
		}
		name, err := filepath.Rel(path.bucketPrefix, *object.Key)
		if err != nil {
			sendErrorInfoToChannel(c, err)
			continue
		}
		if name == "." {
			// Single file was specified
			c <- &fileInfo{
				name:         filepath.Base(*object.Key),
				path:         filepath.Dir(*object.Key),
				size:         *object.Size,
				lastModified: *object.LastModified,
				singleFile:   true,
			}
		} else {
			c <- &fileInfo{
				name:         name,
				path:         *object.Key,
				size:         *object.Size,
				lastModified: *object.LastModified,
			}
		}
	}

	return list.NextContinuationToken
}

// listLocalFiles returns a channel which receives the infos of the files under the given basePath.
// basePath have to be absolute path.
func listLocalFiles(basePath string) chan *fileInfo {
	c := make(chan *fileInfo)

	basePath = filepath.ToSlash(basePath)

	go func() {
		defer close(c)

		stat, err := os.Stat(basePath)
		if os.IsNotExist(err) {
			// The path doesn't exist.
			// Returns and closes the channel without sending any.
			return
		} else if err != nil {
			sendErrorInfoToChannel(c, err)
			return
		}

		if !stat.IsDir() {
			sendFileInfoToChannel(c, filepath.Dir(basePath), basePath, stat, true)
			return
		}

		sendFileInfoToChannel(c, basePath, basePath, stat, false)

		err = filepath.Walk(basePath, func(path string, stat os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			sendFileInfoToChannel(c, basePath, path, stat, false)
			return nil
		})

		if err != nil {
			sendErrorInfoToChannel(c, err)
		}

	}()
	return c
}

func sendFileInfoToChannel(c chan *fileInfo, basePath, path string, stat os.FileInfo, singleFile bool) {
	if stat == nil || stat.IsDir() {
		return
	}
	relPath, _ := filepath.Rel(basePath, path)
	c <- &fileInfo{
		name:         relPath,
		path:         path,
		size:         stat.Size(),
		lastModified: stat.ModTime(),
		singleFile:   singleFile,
	}
}

func sendErrorInfoToChannel(c chan *fileInfo, err error) {
	c <- &fileInfo{
		err: err,
	}
}

// filterFilesForSync filters the source files from the given destination files, and returns
// another channel which includes the files necessary to be synced.
func filterFilesForSync(sourceFileChan, destFileChan chan *fileInfo, del bool) chan *fileOp {
	c := make(chan *fileOp)

	destFiles, err := fileInfoChanToMap(destFileChan)

	go func() {
		defer close(c)
		if err != nil {
			c <- &fileOp{fileInfo: &fileInfo{err: err}}
			return
		}
		for sourceInfo := range sourceFileChan {
			destInfo, ok := destFiles[sourceInfo.name]
			// source is necessary to sync if
			// 1. The dest doesn't exist
			// 2. The dest doesn't have the same size as the source
			// 3. The dest is older than the source
			if !ok || sourceInfo.size != destInfo.size || sourceInfo.lastModified.After(destInfo.lastModified) {
				c <- &fileOp{fileInfo: sourceInfo}
			}
			if ok {
				destInfo.existsInSource = true
			}
		}
		if del {
			for _, destInfo := range destFiles {
				if !destInfo.existsInSource {
					// The source doesn't exist
					c <- &fileOp{fileInfo: destInfo, op: opDelete}
				}
			}
		}
	}()

	return c
}

// fileInfoChanToMap accumulates the fileInfos from the given channel and returns a map.
// It retruns an error if the channel contains an error.
func fileInfoChanToMap(files chan *fileInfo) (map[string]*fileInfo, error) {
	result := make(map[string]*fileInfo)

	for file := range files {
		if file.err != nil {
			return nil, file.err
		}
		result[file.name] = file
	}
	return result, nil
}
