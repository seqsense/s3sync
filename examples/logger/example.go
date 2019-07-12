package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
)

type Logger struct {
}

func (l *Logger) Log(v ...interface{}) {
	fmt.Println(v...)
}

func (l *Logger) Logf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Usage: go run ./examples/simple s3://example-bucket/path/to/source path/to/dest
func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("from=%s\n", os.Args[1])
	fmt.Printf("to=%s\n", os.Args[2])

	s3sync.SetLogger(&Logger{})

	err = s3sync.New(sess).Sync(os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
}
