package shared

import (
	"context"
	"io"
)

// Writer defines the operations for modifying and creating items in the filesystem.
type Writer interface {
	// WriteFile uploads or updates a file with the content from the provided reader.
	WriteFile(ctx context.Context, path string, r io.Reader, opts WriteOptions) (Item, error)
	// Mkdir creates a new directory at the specified path.
	Mkdir(ctx context.Context, path string) error
	// Touch creates a new empty file or updates the modification time of an existing one.
	Touch(ctx context.Context, path string) (Item, error)
}
