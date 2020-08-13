// Copyright 2020 SEQSENSE, Inc.
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
	"strings"
	"sync"
)

type multiErr struct {
	mu  sync.Mutex
	err []error
}

func (e *multiErr) Append(err error) {
	e.mu.Lock()
	e.err = append(e.err, err)
	e.mu.Unlock()
}

func (e *multiErr) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.err)
}

func (e *multiErr) ErrOrNil() error {
	if e.Len() > 0 {
		return e
	}
	return nil
}

func (e *multiErr) Error() string {
	var errMsgs []string
	for _, err := range e.err {
		errMsgs = append(errMsgs, err.Error())
	}
	return strings.Join(errMsgs, "\n")
}
