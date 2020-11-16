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

import "github.com/aws/aws-sdk-go/service/s3/s3manager"

const (
	// Default number of parallel file sync jobs.
	DefaultParallel = 16
)

// Option is a functional option type of Manager.
type Option func(*Manager)

// WithParallel sets maximum number of parallel file sync jobs.
func WithParallel(n int) Option {
	return func(m *Manager) {
		m.nJobs = n
	}
}

// WithDelete enables to delete files unexisting on source directory.
func WithDelete() Option {
	return func(m *Manager) {
		m.del = true
	}
}

// WithACL sets Access Control List string for uploading.
func WithACL(acl string) Option {
	return func(m *Manager) {
		acl := acl
		m.acl = &acl
	}
}

// WithDryRun enables dry-run mode.
func WithDryRun() Option {
	return func(m *Manager) {
		m.dryrun = true
	}
}

// WithoutGuessMimeType disables guessing MIME type from contents.
func WithoutGuessMimeType() Option {
	return func(m *Manager) {
		m.guessMime = false
	}
}

// WithContentType overwrites uploading MIME type.
func WithContentType(mime string) Option {
	return func(m *Manager) {
		m.contentType = &mime
	}
}

// WithDownloaderOptions sets underlying s3manager's options.
func WithDownloaderOptions(opts ...func(*s3manager.Downloader)) Option {
	return func(m *Manager) {
		m.downloaderOpts = opts
	}
}

// WithUploaderOptions sets underlying s3manager's options.
func WithUploaderOptions(opts ...func(*s3manager.Uploader)) Option {
	return func(m *Manager) {
		m.uploaderOpts = opts
	}
}
