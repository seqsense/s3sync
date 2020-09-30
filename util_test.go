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
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getSession() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String("ap-northeast-1"),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4572"),
	})
	return sess
}

type s3Object struct {
	path string
	size int
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

func listObjectsSorted(t *testing.T, bucket string) []s3Object {
	svc := s3.New(session.New(&aws.Config{
		Region:           aws.String("test"),
		Endpoint:         aws.String("http://localhost:4572"),
		S3ForcePathStyle: aws.Bool(true),
	}))

	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  &bucket,
		MaxKeys: aws.Int64(100),
	})
	if err != nil {
		t.Fatal("ListObjects failed", err)
	}
	var objs []s3Object
	for _, obj := range result.Contents {
		objs = append(objs, s3Object{
			path: *obj.Key,
			size: int(*obj.Size),
		})
	}
	sort.Sort(s3ObjectList(objs))
	return objs
}
