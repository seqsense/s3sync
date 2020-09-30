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
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync/atomic"
	"testing"
)

const dummyFilename = "README.md"

func TestS3syncNotImplemented(t *testing.T) {
	m := New(getSession())

	if err := m.Sync("foo", "bar"); err == nil {
		t.Fatal("local to local sync is not supported")
	}

	if err := m.Sync("s3://foo", "s3://bar"); err == nil {
		t.Fatal("s3 to s3 sync is not implemented yet")
	}
}

func TestS3sync(t *testing.T) {
	data, err := ioutil.ReadFile(dummyFilename)
	if err != nil {
		t.Fatal("Failed to read", dummyFilename)
	}

	dummyFileSize := len(data)

	t.Run("Download", func(t *testing.T) {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		destOnlyFilename := filepath.Join(temp, "dest_only_file")
		const destOnlyFileSize = 10
		if err := ioutil.WriteFile(destOnlyFilename, make([]byte, destOnlyFileSize), 0644); err != nil {
			t.Fatal("Failed to write", err)
		}

		// The dummy s3 bucket has following files.
		//
		// s3://example-bucket/
		// ├── README.md
		// ├── bar
		// │   └── baz
		// │       └── README.md
		// └── foo
		//     └── README.md
		if err := New(getSession()).Sync("s3://example-bucket", temp); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		fileHasSize(t, destOnlyFilename, destOnlyFileSize)
		fileHasSize(t, filepath.Join(temp, dummyFilename), dummyFileSize)
		fileHasSize(t, filepath.Join(temp, "foo", dummyFilename), dummyFileSize)
		fileHasSize(t, filepath.Join(temp, "bar/baz", dummyFilename), dummyFileSize)
	})
	t.Run("DownloadSkipDirectory", func(t *testing.T) {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		if err := New(getSession()).Sync("s3://example-bucket-directory", temp); err != nil {
			t.Fatal("Sync should be successful", err)
		}
	})
	t.Run("Upload", func(t *testing.T) {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		for _, dir := range []string{
			filepath.Join(temp, "foo"), filepath.Join(temp, "bar", "baz"),
		} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal("Failed to mkdir", err)
			}
		}

		for _, file := range []string{
			filepath.Join(temp, dummyFilename),
			filepath.Join(temp, "foo", dummyFilename),
			filepath.Join(temp, "bar", "baz", dummyFilename),
		} {
			if err := ioutil.WriteFile(file, make([]byte, dummyFileSize), 0644); err != nil {
				t.Fatal("Failed to write", err)
			}
		}

		if err := New(getSession()).Sync(temp, "s3://example-bucket-upload"); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		objs := listObjectsSorted(t, "example-bucket-upload")
		if n := len(objs); n != 4 {
			t.Fatalf("Number of the files should be 4 (result: %v)", objs)
		}
		for _, obj := range objs {
			if obj.size != dummyFileSize {
				t.Errorf("Object size should be %d, actual %d", dummyFileSize, obj.size)
			}
		}
		if objs[0].path != "README.md" ||
			objs[1].path != "bar/baz/README.md" ||
			objs[2].path != "dest_only_file" ||
			objs[3].path != "foo/README.md" {
			t.Error("Unexpected keys", objs)
		}
	})
}

func TestDelete(t *testing.T) {
	data, err := ioutil.ReadFile(dummyFilename)
	if err != nil {
		t.Fatal("Failed to read", dummyFilename)
	}

	dummyFileSize := len(data)

	t.Run("DeleteLocal", func(t *testing.T) {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		destOnlyFilename := filepath.Join(temp, "dest_only_file")
		const destOnlyFileSize = 10
		if err := ioutil.WriteFile(destOnlyFilename, make([]byte, destOnlyFileSize), 0644); err != nil {
			t.Fatal("Failed to write", err)
		}

		if err := New(getSession(), WithDelete()).Sync(
			"s3://example-bucket", temp,
		); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		if _, err := os.Stat(destOnlyFilename); !os.IsNotExist(err) {
			t.Error("Destination-only-file should be removed by sync")
		}

		fileHasSize(t, filepath.Join(temp, dummyFilename), dummyFileSize)
		fileHasSize(t, filepath.Join(temp, "foo", dummyFilename), dummyFileSize)
		fileHasSize(t, filepath.Join(temp, "bar/baz", dummyFilename), dummyFileSize)
	})
	t.Run("DeleteRemote", func(t *testing.T) {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		for _, dir := range []string{
			filepath.Join(temp, "foo"), filepath.Join(temp, "bar", "baz"),
		} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal("Failed to mkdir", err)
			}
		}

		for _, file := range []string{
			filepath.Join(temp, dummyFilename),
			filepath.Join(temp, "foo", dummyFilename),
			filepath.Join(temp, "bar", "baz", dummyFilename),
		} {
			if err := ioutil.WriteFile(file, make([]byte, dummyFileSize), 0644); err != nil {
				t.Fatal("Failed to write", err)
			}
		}

		if err := New(
			getSession(), WithDelete(),
		).Sync(temp, "s3://example-bucket-delete"); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		objs := listObjectsSorted(t, "example-bucket-delete")
		if n := len(objs); n != 3 {
			t.Fatalf("Number of the files should be 3 (result: %v)", objs)
		}
		for _, obj := range objs {
			if obj.size != dummyFileSize {
				t.Errorf("Object size should be %d, actual %d", dummyFileSize, obj.size)
			}
		}
		if objs[0].path != "README.md" ||
			objs[1].path != "bar/baz/README.md" ||
			objs[2].path != "foo/README.md" {
			t.Error("Unexpected keys", objs)
		}
	})
}

func TestDryRun(t *testing.T) {
	data, err := ioutil.ReadFile(dummyFilename)
	if err != nil {
		t.Fatal("Failed to read", dummyFilename)
	}

	dummyFileSize := len(data)

	t.Run("Download", func(t *testing.T) {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		destOnlyFilename := filepath.Join(temp, "dest_only_file")
		const destOnlyFileSize = 10
		if err := ioutil.WriteFile(destOnlyFilename, make([]byte, destOnlyFileSize), 0644); err != nil {
			t.Fatal("Failed to write", err)
		}

		if err := New(getSession(), WithDelete(), WithDryRun()).Sync(
			"s3://example-bucket", temp,
		); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		fileHasSize(t, destOnlyFilename, destOnlyFileSize)

		if _, err := os.Stat(filepath.Join(temp, dummyFilename)); !os.IsNotExist(err) {
			t.Error("File must not be downloaded on dry-run")
		}
		if _, err := os.Stat(filepath.Join(temp, "foo", dummyFilename)); !os.IsNotExist(err) {
			t.Error("File must not be downloaded on dry-run")
		}
		if _, err := os.Stat(filepath.Join(temp, "bar/baz", dummyFilename)); !os.IsNotExist(err) {
			t.Error("File must not be downloaded on dry-run")
		}
	})
	t.Run("Upload", func(t *testing.T) {
		temp, err := ioutil.TempDir("", "s3synctest")
		defer os.RemoveAll(temp)

		if err != nil {
			t.Fatal("Failed to create temp dir")
		}

		for _, dir := range []string{
			filepath.Join(temp, "foo"), filepath.Join(temp, "bar", "baz"),
		} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal("Failed to mkdir", err)
			}
		}

		for _, file := range []string{
			filepath.Join(temp, dummyFilename),
			filepath.Join(temp, "foo", dummyFilename),
			filepath.Join(temp, "bar", "baz", dummyFilename),
		} {
			if err := ioutil.WriteFile(file, make([]byte, dummyFileSize), 0644); err != nil {
				t.Fatal("Failed to write", err)
			}
		}

		if err := New(
			getSession(), WithDelete(), WithDryRun(),
		).Sync(temp, "s3://example-bucket-dryrun"); err != nil {
			t.Fatal("Sync should be successful", err)
		}

		objs := listObjectsSorted(t, "example-bucket-delete")
		if n := len(objs); n != 1 {
			t.Fatalf("Number of the files should be 1 (result: %v)", objs)
		}
		if n := objs[0].size; n != dummyFileSize {
			t.Errorf("Object size should be %d, actual %d", dummyFileSize, n)
		}
		if objs[0].path != "dest_only_file" {
			t.Error("Unexpected key", objs[0].path)
		}
	})
}

func TestPartialS3sync(t *testing.T) {
	data, err := ioutil.ReadFile(dummyFilename)
	if err != nil {
		t.Fatal("Failed to read", dummyFilename)
	}

	expectedFileSize := len(data)

	temp, err := ioutil.TempDir("", "s3synctest")
	defer os.RemoveAll(temp)

	if err != nil {
		t.Fatal("Failed to create temp dir")
	}

	var syncCount uint32
	SetLogger(createLoggerWithLogFunc(func(v ...interface{}) {
		atomic.AddUint32(&syncCount, 1) // This function is called once per one download
	}))

	if err := New(getSession()).Sync("s3://example-bucket", temp); err != nil {
		t.Fatal("Sync should be successful", err)
	}

	if atomic.LoadUint32(&syncCount) != 3 {
		t.Fatal("3 files should be synced")
	}

	atomic.StoreUint32(&syncCount, 0)

	os.RemoveAll(filepath.Join(temp, "foo"))

	if New(getSession()).Sync("s3://example-bucket", temp) != nil {
		t.Fatal("Sync should be successful")
	}

	if atomic.LoadUint32(&syncCount) != 1 {
		t.Fatal("Only 1 file should be synced")
	}

	fileHasSize(t, filepath.Join(temp, dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "foo", dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "bar/baz", dummyFilename), expectedFileSize)
	// TODO: Assert only one file was downloaded at the second sync.
}

func TestListLocalFiles(t *testing.T) {
	temp, err := ioutil.TempDir("", "s3synctest")
	defer os.RemoveAll(temp)

	if err != nil {
		t.Fatal("Failed to create temp dir")
	}

	for _, dir := range []string{
		filepath.Join(temp, "empty"),
		filepath.Join(temp, "foo"),
		filepath.Join(temp, "bar", "baz"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal("Failed to mkdir", err)
		}
	}

	for _, file := range []string{
		filepath.Join(temp, "test1"),
		filepath.Join(temp, "foo", "test2"),
		filepath.Join(temp, "bar", "baz", "test3"),
	} {
		if err := ioutil.WriteFile(file, make([]byte, 10), 0644); err != nil {
			t.Fatal("Failed to write", err)
		}
	}

	collectFilePaths := func(ch chan *fileInfo) []string {
		list := []string{}
		for f := range ch {
			list = append(list, f.path)
		}
		sort.Strings(list)
		return list
	}

	t.Run("Root", func(t *testing.T) {
		paths := collectFilePaths(listLocalFiles(temp))
		expected := []string{
			filepath.Join(temp, "bar", "baz", "test3"),
			filepath.Join(temp, "foo", "test2"),
			filepath.Join(temp, "test1"),
		}
		if !reflect.DeepEqual(expected, paths) {
			t.Errorf("Local file list is expected to be %v, got %v", expected, paths)
		}
	})

	t.Run("EmptyDir", func(t *testing.T) {
		paths := collectFilePaths(listLocalFiles(filepath.Join(temp, "empty")))
		expected := []string{}
		if !reflect.DeepEqual(expected, paths) {
			t.Errorf("Local file list is expected to be %v, got %v", expected, paths)
		}
	})

	t.Run("File", func(t *testing.T) {
		paths := collectFilePaths(listLocalFiles(filepath.Join(temp, "test1")))
		expected := []string{
			filepath.Join(temp, "test1"),
		}
		if !reflect.DeepEqual(expected, paths) {
			t.Errorf("Local file list is expected to be %v, got %v", expected, paths)
		}
	})

	t.Run("Dir", func(t *testing.T) {
		paths := collectFilePaths(listLocalFiles(filepath.Join(temp, "foo")))
		expected := []string{
			filepath.Join(temp, "foo", "test2"),
		}
		if !reflect.DeepEqual(expected, paths) {
			t.Errorf("Local file list is expected to be %v, got %v", expected, paths)
		}
	})

	t.Run("Dir2", func(t *testing.T) {
		paths := collectFilePaths(listLocalFiles(filepath.Join(temp, "bar")))
		expected := []string{
			filepath.Join(temp, "bar", "baz", "test3"),
		}
		if !reflect.DeepEqual(expected, paths) {
			t.Errorf("Local file list is expected to be %v, got %v", expected, paths)
		}
	})
}

type dummyLogger struct {
	log func(...interface{})
}

func (d *dummyLogger) Log(v ...interface{}) {
	d.log(v...)
}
func (d *dummyLogger) Logf(format string, v ...interface{}) {
}

func createLoggerWithLogFunc(log func(v ...interface{})) LoggerIF {
	return &dummyLogger{log: log}
}

func fileHasSize(t *testing.T, filename string, expectedSize int) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(filename, "is not synced")
	}
	if len(data) != expectedSize {
		t.Fatal(filename, "is not synced")
	}
}
