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
	"fmt"
	"net/url"
	"testing"
)

func TestURLToS3Path(t *testing.T) {
	t.Run("NoBucketName", func(t *testing.T) {
		_, err := urlToS3Path(&url.URL{
			Host: "",
			Path: "test",
		})
		if err != errNoBucketName {
			t.Fatalf("Expected error %v, got %v", errNoBucketName, err)
		}
	})
	t.Run("Normal", func(t *testing.T) {
		p, err := urlToS3Path(&url.URL{
			Host: "bucket",
			Path: "test",
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if p.bucket != "bucket" || p.bucketPrefix != "test" {
			t.Fatalf("Unexpected s3Path: %v", p)
		}
	})
	t.Run("UrlEscapedPath", func(t *testing.T) {
		urlHost := "bucket"
		urlPath := url.QueryEscape("test/it")
		srUrl, err := url.Parse(fmt.Sprintf("s3://%s/%s", urlHost, urlPath))

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		p, err := urlToS3Path(srUrl)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if p.bucket != urlHost || p.bucketPrefix != urlPath {
			t.Fatalf("Unexpected s3Path: %v", p)
		}
	})
}

func TestS3Path_String(t *testing.T) {
	p := &s3Path{
		bucket:       "bucket",
		bucketPrefix: "test",
	}
	const expected = "s3://bucket/test"
	if s := p.String(); s != expected {
		t.Fatalf("Expected %s, got %s", expected, s)
	}
}
