# Migration guide

## v2

- aws/aws-sdk-go is updated to aws/aws-sdk-go-v2
  - 🔄Create s3sync client with `aws.Config`
    ```diff
     import (
    -  "github.com/aws/aws-sdk-go/aws/session"
    -  "github.com/seqsense/s3sync"
    +  "github.com/aws/aws-sdk-go-v2/config"
    +  "github.com/seqsense/s3sync/v2"
     )

    -sess, err := session.NewSession()
    +cfg, err := config.LoadDefaultConfig(ctx)
     if err != nil {
       // error handling
     }
    -syncManager := s3sync.New(sess)
    +syncManager := s3sync.New(cfg)
    ```
- `Sync()` requires context by default and `SyncWithContext()` is removed
  - 🔄Pass ctx to `Sync()`
    ```diff
    -syncManager.Sync("s3://bucket/key", "local/path")
    +syncManager.Sync(ctx, "s3://bucket/key", "local/path")
    ```
  - 🔄Use `Sync()` instead of `SyncWithContext()`
    ```diff
    -syncManager.SyncWithContext(ctx, "s3://bucket/key", "local/path")
    +syncManager.Sync(ctx, "s3://bucket/key", "local/path")
    ```
