package shared

import (
	"context"
	"io"
)

// Writer defines the operations for modifying and creating items in the filesystem.
type Writer interface {
	// WriteFile uploads or updates a file with the content from the provided reader.
	WriteFile(ctx context.Context, path string, r io.Reader) (Item, error)
	// Mkdir creates a new directory at the specified path.
	Mkdir(ctx context.Context, path string) error
}
