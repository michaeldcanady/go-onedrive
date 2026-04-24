package onedrive

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

const (
	// uploadChunkSize is the size of each chunk in a resumable upload (must be multiple of 320 KiB).
	uploadChunkSize = 320 * 1024 * 10 // 3.2 MiB
)

func writeLargeFile(ctx context.Context, b *Backend, token, driveID, itemPath string, r io.Reader, opts fs.WriteOptions) (fs.Item, error) {
	return fs.Item{}, fmt.Errorf("writeLargeFile not implemented in decoupled backend")
}
