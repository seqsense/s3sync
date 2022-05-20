# s3sync

![CI](https://github.com/seqsense/s3sync/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/seqsense/s3sync/branch/master/graph/badge.svg)](https://codecov.io/gh/seqsense/s3sync)

> Golang utility for syncing between s3 and local

# What's new

## v2

- Default AWS SDK is updated to aws-sdk-go-v2.
  You can still use aws-sdk-go (v1) [by this code](#using-aws-sdk-go-v1).
- `Sync()` receives `context.Context` as a first argument.

# Usage

Use `New` to create a manager, and `Sync` function syncs between s3 and local filesystem.

```go
import (
  "context"

  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/seqsense/s3sync"
)

func main() {
  // Creates an AWS session
  cfg, _ := config.LoadDefaultConfig(context.TODO())

  syncManager := s3sync.New(cfg)

  // Sync from s3 to local
  syncManager.Sync(context.TODO(), "s3://yourbucket/path/to/dir", "local/path/to/dir")

  // Sync from local to s3
  syncManager.Sync(context.TODO(), "local/path/to/dir", "s3://yourbucket/path/to/dir")
}
```

- Note: Sync from s3 to s3 is not implemented yet.

## Setting the custom logger

You can set your custom logger.

```go
import "github.com/seqsense/s3sync"

...
s3sync.SetLogger(&CustomLogger{})
...
```

The logger needs to implement `Log` and `Logf` methods. See the godoc for details.

## Setting up the parallelism

You can configure the number of parallel jobs for sync. Default is 16.

```
s3sync.new(sess, s3sync.WithParallel(16)) // This is the same as default.
s3sync.new(sess, s3sync.WithParallel(1)) // You can sync one by one.
```

## Using aws-sdk-go v1

```go
import (
  "github.com/at-wat/s3iot/awssdkv1"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3manager"
)

func main() {
  // Creates an AWS session
  sess, _ := session.NewSession(&aws.Config{
    Region: aws.String("us-east-1"),
  })
  cli := s3.New(sess)

  // Initialize s3sync.Manager using aws-sdk-go (v1)
  syncManager := s3sync.NewFromAPI(s3iotiface.CombineUpDownloader(
    awssdkv1.NewAWSSDKUploader(s3manager.NewUploaderWithClient(cli)),
    awssdkv1.NewAWSSDKDownloader(s3manager.NewDownloaderWithClient(cli)),
  ), awssdkv1.NewAPI(cli))

  // ...
```

# License

Apache 2.0 License. See [LICENSE](https://github.com/seqsense/s3sync/blob/master/LICENSE).
