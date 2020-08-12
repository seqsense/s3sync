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
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const dummyFilename = "README.md"

func TestS3syncNotImplemented(t *testing.T) {
	m := New(getSession())

	if err := m.Sync("foo", "bar"); err == nil {
		t.Fatal("local to local sync is not supported")
	}

	if err := m.Sync("s3://foo", "s3://bar"); err == nil {
		t.Fatal("s3 to s3 sync is not implemented yet")
	}
}

func TestS3sync(t *testing.T) {
	t.Run("Download", func(t *testing.T) {
		data, err := ioutil.ReadFile(dummyFilename)
		if err != nil {
			t.Fatal("Failed to read", dummyFilename)
		}

		expectedFileSize := len(data)

		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		// The dummy s3 bucket has following files.
		//
		// s3://example-bucket/
		// ├── README.md
		// ├── bar
		// │   └── baz
		// │       └── README.md
		// └── foo
		//     └── README.md
		if err := New(getSession()).Sync("s3://example-bucket", temp); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		fileHasSize(t, filepath.Join(temp, dummyFilename), expectedFileSize)
		fileHasSize(t, filepath.Join(temp, "foo", dummyFilename), expectedFileSize)
		fileHasSize(t, filepath.Join(temp, "bar/baz", dummyFilename), expectedFileSize)
	})
	t.Run("Upload", func(t *testing.T) {
		const expectedFileSize = 100

		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		for _, dir := range []string{
			filepath.Join(temp, "foo"), filepath.Join(temp, "bar", "baz"),
		} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal("Failed to mkdir", err)
			}
		}

		for _, file := range []string{
			filepath.Join(temp, dummyFilename),
			filepath.Join(temp, "foo", dummyFilename),
			filepath.Join(temp, "bar", "baz", dummyFilename),
		} {
			if err := ioutil.WriteFile(file, make([]byte, expectedFileSize), 0644); err != nil {
				t.Fatal("Failed to write", err)
			}
		}

		if err := New(getSession()).Sync(temp, "s3://example-bucket-upload"); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		svc := s3.New(session.New(&aws.Config{
			Region:           aws.String("test"),
			Endpoint:         aws.String("http://localhost:4572"),
			S3ForcePathStyle: aws.Bool(true),
		}))

		result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:  aws.String("example-bucket-upload"),
			MaxKeys: aws.Int64(10),
		})
		if err != nil {
			t.Fatal("ListObjects failed", err)
		}
		if n := len(result.Contents); n != 3 {
			t.Fatalf("Number of the files should be 3 (result: %s)", result)
		}
		var keys []string
		for _, obj := range result.Contents {
			if int(*obj.Size) != expectedFileSize {
				t.Errorf("Object size should be %d, actual %d", expectedFileSize, obj.Size)
			}
			keys = append(keys, *obj.Key)
		}
		sort.Strings(keys)
		if keys[0] != "README.md" ||
			keys[1] != "bar/baz/README.md" ||
			keys[2] != "foo/README.md" {
			t.Error("Unexpected keys", keys)
		}
	})
}

func TestPartialS3sync(t *testing.T) {
	data, err := ioutil.ReadFile(dummyFilename)
	if err != nil {
		t.Fatal("Failed to read", dummyFilename)
	}

	expectedFileSize := len(data)

	temp, err := ioutil.TempDir("", "s3synctest")
	defer os.RemoveAll(temp)

	if err != nil {
		t.Fatal("Failed to create temp dir")
	}

	syncCount := 0
	SetLogger(createLoggerWithLogFunc(func(v ...interface{}) {
		syncCount++ // This function is called once per one download
	}))

	if err := New(getSession()).Sync("s3://example-bucket", temp); err != nil {
		t.Fatal("Sync should be successful", err)
	}

	if syncCount != 3 {
		t.Fatal("3 files should be synced")
	}

	syncCount = 0

	os.RemoveAll(filepath.Join(temp, "foo"))

	if New(getSession()).Sync("s3://example-bucket", temp) != nil {
		t.Fatal("Sync should be successful")
	}

	if syncCount != 1 {
		t.Fatal("Only 1 file should be synced")
	}

	fileHasSize(t, filepath.Join(temp, dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "foo", dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "bar/baz", dummyFilename), expectedFileSize)
	// TODO: Assert only one file was downloaded at the second sync.
}

func getSession() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String("ap-northeast-1"),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4572"),
	})
	return sess
}

type dummyLogger struct {
	log func(...interface{})
}

func (d *dummyLogger) Log(v ...interface{}) {
	d.log(v...)
}
func (d *dummyLogger) Logf(format string, v ...interface{}) {
}

func createLoggerWithLogFunc(log func(v ...interface{})) LoggerIF {
	return &dummyLogger{log: log}
}

func fileHasSize(t *testing.T, filename string, expectedSize int) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(filename, "is not synced")
	}
	if len(data) != expectedSize {
		t.Fatal(filename, "is not synced")
	}
}
