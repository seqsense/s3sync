# s3sync

> An utility for golang which syncs between s3 and local disk

# Usage

```go
import "github.com/seqsense/s3sync"

func main() {
  // Sync from s3 to local
  s3sync.Sync("s3://yourbucket/path/to/dir", "local/path/to/dir")

  // Sync from local to s3
  s3sync.Sync("local/path/to/dir", "s3://yourbucket/path/to/dir")

  // Sync from s3 to s3
  s3sync.Sync("s3://yourbucket/path/to/dir", "s3://anotheryourbucket/path/to/dir")
}
```
