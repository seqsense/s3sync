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
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
)

// Usage: go run ./examples/simple s3://example-bucket/path/to/source path/to/dest
func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("from=%s\n", os.Args[1])
	fmt.Printf("to=%s\n", os.Args[2])

	manager := s3sync.New(sess)
	err = manager.Sync(os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
	s := manager.GetStatistics()
	fmt.Printf("Sync results:\nBytes written: %d\nFiles uploaded: %d\nTime spent: %d\nFiles deleted: %d\n", s.Bytes(), s.SyncTime(), s.Files(), s.DeletedFiles())
}
