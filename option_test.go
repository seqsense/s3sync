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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestWithParallel(t *testing.T) {
	m := New(aws.Config{}, WithParallel(2))
	if m.nJobs != 2 {
		t.Fatal("Manager.nJobs must be configured by WithParallel option")
	}
}

func TestWithACL(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		m := New(aws.Config{})
		if m.acl != nil {
			t.Fatal("Manager.acl must be nil if initialized with WithACL")
		}
	})
	t.Run("WithACL", func(t *testing.T) {
		m := New(aws.Config{}, WithACL("test"))
		if *m.acl != "test" {
			t.Fatal("Manager.acl must be configured by WithParallel option")
		}
	})
}
