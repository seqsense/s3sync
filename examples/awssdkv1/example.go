// Copyright 2022 SEQSENSE, Inc.
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

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/at-wat/s3iot/awss3v1"
	"github.com/at-wat/s3iot/s3iotiface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/seqsense/s3sync"
)

// Usage: go run ./examples/simple s3://example-bucket/path/to/source path/to/dest
func main() {
	// Creates an AWS session
	sess, err := session.NewSession(&aws.Config{
		// Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}
	cli := s3.New(sess)

	// Initialize s3sync.Manager using aws-sdk-go (v1)
	syncManager := s3sync.NewFromAPI(s3iotiface.CombineUpDownloader(
		awss3v1.NewAWSSDKUploader(s3manager.NewUploaderWithClient(cli)),
		awss3v1.NewAWSSDKDownloader(s3manager.NewDownloaderWithClient(cli)),
	), awss3v1.NewAPI(cli))

	fmt.Printf("from=%s\n", os.Args[1])
	fmt.Printf("to=%s\n", os.Args[2])

	err = syncManager.Sync(context.TODO(), os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
}
