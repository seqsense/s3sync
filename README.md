# s3sync

![CI](https://github.com/seqsense/s3sync/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/seqsense/s3sync/branch/master/graph/badge.svg)](https://codecov.io/gh/seqsense/s3sync)

> Golang utility for syncing between s3 and local

## Migration guide

- [v2](MIGRATION.md#v2)

# Usage

Use `New` to create a manager, and `Sync` function syncs between s3 and local filesystem.

```go
import (
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/seqsense/s3sync/v2"
)

func main() {
  ctx := context.TODO()

  // Load AWS config
  cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
  if err != nil {
    panic(err)
  }

  syncManager := s3sync.New(cfg)

  // Sync from s3 to local
  syncManager.Sync(ctx, "s3://yourbucket/path/to/dir", "local/path/to/dir")

  // Sync from local to s3
  syncManager.Sync(ctx, "local/path/to/dir", "s3://yourbucket/path/to/dir")

  // Sync from s3 to s3
  syncManager.Sync(ctx, "s3://yourbucket/path/to/dir", "s3://anotherbucket/path/to/dir")
}
```

## Sets the custom logger

You can set your custom logger.

```go
import "github.com/seqsense/s3sync/v2"

...
s3sync.SetLogger(&CustomLogger{})
...
```

The logger needs to implement `Logf` methods. See the godoc for details.

## Sets up the parallelism

You can configure the number of parallel jobs for sync. Default is 16.

```
s3sync.New(cfg, s3sync.WithParallel(16)) // This is the same as default.
s3sync.New(cfg, s3sync.WithParallel(1)) // You can sync one by one.
```

# License

Apache 2.0 License. See [LICENSE](https://github.com/seqsense/s3sync/blob/master/LICENSE).
