# Migration guide

## v3

- aws-sdk-go-v2/feature/s3/manager is replaced by aws-sdk-go-v2/feature/s3/transfermanager
  - 🔄Pass `transfermanager.Options` modifiers to `WithUploaderOptions()` and `WithDownloaderOptions()`
    ```diff
     import (
    -  "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    -  "github.com/seqsense/s3sync/v2"
    +  "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
    +  "github.com/seqsense/s3sync/v3"
     )

     syncManager := s3sync.New(cfg,
    -  s3sync.WithUploaderOptions(func(u *manager.Uploader) {
    -    u.Concurrency = 3
    -    u.PartSize = 16 * 1024 * 1024
    +  s3sync.WithUploaderOptions(func(o *transfermanager.Options) {
    +    o.Concurrency = 3
    +    o.PartSizeBytes = 16 * 1024 * 1024
       }),
     )
    ```
- Go 1.24 or later is required

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
