package s3sync

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const dummyFilename = "README.md"

func TestS3syncNotImplemented(t *testing.T) {
	m := New(getSession())

	if err := m.Sync("foo", "bar"); err == nil {
		t.Fatal("local to local sync is not supported")
	}

	if err := m.Sync("foo", "s3://bar"); err == nil {
		t.Fatal("local to s3 sync is not implemented yet")
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

	expectedFileSize := len(data)

	temp, err := ioutil.TempDir("", "s3synctest")
	defer os.RemoveAll(temp)

	if err != nil {
		t.Fatal("Failed to create temp dir")
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
	if New(getSession()).Sync("s3://example-bucket", temp) != nil {
		t.Fatal("Sync should be successful")
	}

	fileHasSize(t, filepath.Join(temp, dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "foo", dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "bar/baz", dummyFilename), expectedFileSize)
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

	syncCount := 0
	SetLogger(createLoggerWithLogFunc(func(v ...interface{}) {
		syncCount++ // This function is called once per one download
	}))

	if New(getSession()).Sync("s3://example-bucket", temp) != nil {
		t.Fatal("Sync should be successful")
	}

	if syncCount != 3 {
		t.Fatal("3 files should be synced")
	}

	syncCount = 0

	os.RemoveAll(filepath.Join(temp, "foo"))

	if New(getSession()).Sync("s3://example-bucket", temp) != nil {
		t.Fatal("Sync should be successful")
	}

	if syncCount != 1 {
		t.Fatal("Only 1 file should be synced")
	}

	fileHasSize(t, filepath.Join(temp, dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "foo", dummyFilename), expectedFileSize)
	fileHasSize(t, filepath.Join(temp, "bar/baz", dummyFilename), expectedFileSize)
	// TODO: Assert only one file was downloaded at the second sync.
}

func getSession() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String("ap-northeast-1"),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4572"),
	})
	return sess
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
