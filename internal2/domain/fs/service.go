// internal/fs/service.go
package fs

import (
	"context"
	"io"
)

type Service interface {
	Get(ctx context.Context, path string) (Item, error)
	List(ctx context.Context, path string, opts ListOptions) ([]Item, error)
	Stat(ctx context.Context, path string, opts StatOptions) (Item, error)
	ReadFile(ctx context.Context, path string, opts ReadOptions) (io.ReadCloser, error)
	WriteFile(ctx context.Context, path string, r io.Reader, opts WriteOptions) (Item, error)
	Mkdir(ctx context.Context, path string, opts MKDirOptions) error
	Remove(ctx context.Context, path string, opts RemoveOptions) error
	Move(ctx context.Context, src, dst string, opts MoveOptions) error
	Upload(ctx context.Context, src, dst string, opts UploadOptions) (Item, error)
}
