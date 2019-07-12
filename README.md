# s3sync

> Golang utility for syncing between s3 and local

# Usage

```go
import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/seqsense/s3sync"
)

func main() {
  // Creates an AWS session
  sess, _ := session.NewSession(&aws.Config{
    Region: aws.String("us-east-1"),
  })

  syncManager := s3sync.New(sess)
  // Sync from s3 to local
  syncManager.Sync("s3://yourbucket/path/to/dir", "local/path/to/dir")
}
```

- Note: Sync from local to s3 is not implemented yet.
- Note: Sync from s3 to s3 is not implemented yet as well.

## Sets the custom logger

You can set your custom logger.

```go
import "github.com/seqsense/s3sync"

...
s3sync.SetLogger(&CustomLogger{})
...
```

The logger needs to implement `Log` and `Logf` methods. See the godoc for details.
