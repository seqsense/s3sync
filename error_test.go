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
	"errors"
	"testing"
)

func TestMultiErr(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		err := &multiErr{}
		if err.Error() != "" {
			t.Error("Empty multiErr should return empty string")
		}
		if err.ErrOrNil() != nil {
			t.Error("Empty multiErr should return nil error")
		}
	})
	t.Run("MultipleErrors", func(t *testing.T) {
		err := &multiErr{}
		err.Append(errors.New("error1"))
		err.Append(errors.New("error2"))
		if err.Error() != "error1\nerror2" {
			t.Error("Empty multiErr should return joined error message")
		}
		if err.ErrOrNil() != err {
			t.Error("Empty multiErr should return self pointer")
		}
	})
}
