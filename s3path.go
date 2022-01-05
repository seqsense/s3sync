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
	"path/filepath"
	"strings"
)

var errNoBucketName = errors.New("s3 url is missing bucket name")

type s3Path struct {
	bucket       string
	bucketPrefix string
}

func urlToS3Path(url *url.URL) (*s3Path, error) {
	if url.Host == "" {
		return nil, errNoBucketName
	}

	return &s3Path{
		bucket: url.Host,
		// Using filepath.ToSlash for change backslash to slash on Windows
		bucketPrefix: strings.TrimPrefix(filepath.ToSlash(url.Path), "/"),
	}, nil
}

func (p *s3Path) String() string {
	return "s3://" + p.bucket + "/" + p.bucketPrefix
}
