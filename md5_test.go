// Copyright 2021 SEQSENSE, Inc.
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
	"testing"
)

func TestGetMd5Hash(t *testing.T) {
	// The hash value should be obtained from a file
	// with the same content that was created in advance.
	tests := []struct {
		name     string
		content  string
		wantHash string
		hasError bool
	}{
		{
			"files with correct hash values",
			"This is a test sentence.",
			"4d6a0c4cf3f07eadd5ba147c67c6896f",
			false,
		},
		{
			"files with correct hash values",
			"column1,column2,column3\ntest,test,test",
			"1ed7dedaf1bfac6642b52a19a4a1988c",
			false,
		},
		{
			"not existed file",
			"",
			"",
			true,
		},
	}

	for _, tt := range tests {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		t.Run(tt.name, func(t *testing.T) {
			var testFile *os.File
			var testFileName string

			if !tt.hasError {
				testFile, err = ioutil.TempFile(temp, "s3synctest-")
				if err != nil {
					t.Fatal("Failed to create temp file")
				}
				testFileName = testFile.Name()
				defer os.Remove(testFileName)

				if _, err = testFile.Write([]byte(tt.content)); err != nil {
					t.Fatal("Failed to write temp file")
				}
			}

			got, err := getMd5Hash(testFileName)
			if !tt.hasError && err != nil {
				t.Errorf("got err. err:%v, want: err==nil", err)
			}

			if got != tt.wantHash {
				t.Fatalf("Unexpected hash. got: %v, want: %v", got, tt.wantHash)
			}
		})
	}
}
