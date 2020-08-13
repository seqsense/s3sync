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
