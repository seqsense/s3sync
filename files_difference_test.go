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

import "testing"

func TestFilesDifference(t *testing.T) {
	t.Run("false", func(t *testing.T) {
		hasDifference := &filesDifference{}
		if hasDifference.Get() {
			t.Error("filesDifference should return false")
		}
	})
	t.Run("true", func(t *testing.T) {
		hasDifference := &filesDifference{}
		hasDifference.Set(true)
		if !hasDifference.Get() {
			t.Error("filesDifference should return true")
		}
	})
}
