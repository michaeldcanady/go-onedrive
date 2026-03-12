package shared

import (
	"context"
	"io"
)

// Reader defines the operations for retrieving item information and file contents.
type Reader interface {
	// Get retrieves a single Item by its path.
	Get(ctx context.Context, path string) (Item, error)
	// List returns the immediate children of the specified directory path.
	List(ctx context.Context, path string, opts ListOptions) ([]Item, error)
	// ReadFile provides an io.ReadCloser for the content of the file at path.
	ReadFile(ctx context.Context, path string, opts ReadOptions) (io.ReadCloser, error)
	// Stat returns metadata for an item at the specified path.
	Stat(ctx context.Context, path string) (Item, error)
}
