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

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/at-wat/s3iot/awss3v2"
	"github.com/at-wat/s3iot/s3iotiface"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/seqsense/s3sync"
)

// Usage: go run ./examples/simple s3://example-bucket/path/to/source path/to/dest
func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	s3cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Configure S3 options
		o.UsePathStyle = true // Use https://s3.amazonaws.com/BUCKET/KEY instead of https://BUCKET.s3.amazonaws.com/KEY
	})
	s3api := awss3v2.NewAPI(s3cli)
	syncManager := s3sync.NewFromAPI(s3iotiface.CombineUpDownloader(
		awss3v2.NewAWSSDKUploader(manager.NewUploader(s3cli, func(u *manager.Uploader) {
			// Configure uploader options
			u.Concurrency = 2 // Limit concurrency of each file upload
		})),
		awss3v2.NewAWSSDKDownloader(manager.NewDownloader(s3cli, func(u *manager.Downloader) {
			// Configure downloader options
			u.Concurrency = 2 // Limit concurrency of each file download
		})),
	), s3api, s3sync.WithParallel(3))

	// Note: in this configuration, three files are concurrently up/downloaded and each of them is
	//       up/downloaded by two parallel connections
	//       (so, the total maximum concurrent connections will be 6)

	fmt.Printf("from=%s\n", os.Args[1])
	fmt.Printf("to=%s\n", os.Args[2])

	err = syncManager.Sync(context.TODO(), os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
}
