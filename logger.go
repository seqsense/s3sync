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

import "log"

// LoggerIF is the logger interface which this library requires.
type LoggerIF interface {
	// Log inserts a log entry. Arguments are handled in the manner
	// of fmt.Print.
	Log(v ...interface{})
	// Log inserts a log entry. Arguments are handled in the manner
	// of fmt.Printf.
	Logf(format string, v ...interface{})
}

// Logger is the logger instance.
var logger LoggerIF

// SetLogger sets the logger.
func SetLogger(l LoggerIF) {
	logger = l
}

func println(v ...interface{}) {
	if logger == nil {
		log.Println(v...)
		return
	}
	logger.Log(v...)
}
