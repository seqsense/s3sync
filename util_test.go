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
	"context"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getConfig() aws.Config {
	return aws.Config{
		Region: "ap-northeast-1",
		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           "http://localhost:4572",
				PartitionID:   "localstack",
				SigningRegion: region,
				SigningName:   service,
				SigningMethod: "v4",
			}, nil
		}),
		Credentials: credentials.NewStaticCredentialsProvider(
			"ASIAZZZZZZZZZZZZZZZZ",
			"0000000000000000000000000000000000000000",
			"",
		),
	}
}

func usePathStyle(o *s3.Options) {
	o.UsePathStyle = true
}

func newManager(options ...Option) *Manager {
	cli := s3.NewFromConfig(getConfig(), usePathStyle)
	return newFromS3ClientV2(cli, options...)
}

type s3Object struct {
	path        string
	size        int
	contentType string
}

type s3ObjectList []s3Object

func (l s3ObjectList) Len() int {
	return len(l)
}
func (l s3ObjectList) Less(i, j int) bool {
	return l[i].path < l[j].path
}
func (l s3ObjectList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func deleteObject(t *testing.T, bucket, key string) {
	svc := s3.NewFromConfig(getConfig(), usePathStyle)

	_, err := svc.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatal("DeleteObject failed", err)
	}
}

func listObjectsSorted(t *testing.T, bucket string) []s3Object {
	svc := s3.NewFromConfig(getConfig(), usePathStyle)

	result, err := svc.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket:  &bucket,
		MaxKeys: 100,
	})
	if err != nil {
		t.Fatal("ListObjects failed", err)
	}
	var objs []s3Object
	for _, obj := range result.Contents {
		o, err := svc.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    obj.Key,
		})
		if err != nil {
			t.Fatal("GetObject failed", err)
		}
		objs = append(objs, s3Object{
			path:        *obj.Key,
			size:        int(obj.Size),
			contentType: *o.ContentType,
		})
	}
	sort.Sort(s3ObjectList(objs))
	return objs
}

func fileHasSize(t *testing.T, filename string, expectedSize int) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(filename, "is not synced")
		return
	}
	if n := len(data); n != expectedSize {
		t.Errorf("%s is not synced (file size is expected to be %d, actual %d)", filename, expectedSize, n)
	}
}

func fileModTimeBefore(t *testing.T, filename string, t0 time.Time) {
	info, err := os.Stat(filename)
	if err != nil {
		t.Error("Failed to get stat:", err)
		return
	}
	if t1 := info.ModTime(); !t1.Before(t0) {
		t.Errorf("File modification time %v is later than %v", t1, t0)
	}
}
