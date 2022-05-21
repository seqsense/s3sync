# s3sync

![CI](https://github.com/seqsense/s3sync/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/seqsense/s3sync/branch/master/graph/badge.svg)](https://codecov.io/gh/seqsense/s3sync)

> Golang utility for syncing between s3 and local

# What's new

## v2

- `Sync()` receives `context.Context` as a first argument.
- Default AWS SDK is updated to aws-sdk-go-v2.
  See an [example](examples/awssdkv1/example.go) how to use aws-sdk-go (v1) in your code.
- `WithDownloaderOptions()` and `WithUploaderOptions()` are removed.
  See an [example](examples/parallel/example.go) how to set uploader/downloader options.

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

You can globally set your custom logger.

```go
s3sync.SetLogger(&CustomLogger{})
```

The logger needs to implement `Log` and `Logf` methods. See the godoc for details.

[example code](examples/logger/example.go)

## Setting up the parallelism

You can configure the number of parallel jobs for sync. Default is 16.
Note that each file may be transferred in parallel according to the underlying uploader/downloader implementation.

```go
s3sync.New(cfg, s3sync.WithParallel(16)) // This is the same as default.
s3sync.New(cfg, s3sync.WithParallel(1))  // You can sync one by one.
```

[example code](examples/parallel/example.go)

# License

Apache 2.0 License. See [LICENSE](https://github.com/seqsense/s3sync/blob/master/LICENSE).
