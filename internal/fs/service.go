package fs

import (
	"context"
	"io"
)

// Service combines the Reader, Writer, and Manager interfaces to provide a full-featured filesystem service.
type Service interface {
	Namer
	Reader
	Writer
	Manager
}

// Namer provides the name of the filesystem provider.
type Namer interface {
	Name() string
}

// Reader defines the operations for retrieving item information and file contents.
type Reader interface {
	// Get retrieves a single Item by its path.
	Get(ctx context.Context, path *URI) (Item, error)
	// List returns the immediate children of the specified directory path.
	List(ctx context.Context, path *URI, opts ListOptions) ([]Item, error)
	// ReadFile provides an io.ReadCloser for the content of the file at path.
	ReadFile(ctx context.Context, path *URI, opts ReadOptions) (io.ReadCloser, error)
	// Stat returns metadata for an item at the specified path.
	Stat(ctx context.Context, path *URI) (Item, error)
}

// Writer defines the operations for modifying and creating items in the filesystem.
type Writer interface {
	// WriteFile uploads or updates a file with the content from the provided reader.
	WriteFile(ctx context.Context, path *URI, r io.Reader, opts WriteOptions) (Item, error)
	// Mkdir creates a new directory at the specified path.
	Mkdir(ctx context.Context, path *URI) error
	// Touch creates a new empty file or updates the modification time of an existing one.
	Touch(ctx context.Context, path *URI) (Item, error)
}

// Manager defines the operations for higher-level filesystem management and item manipulation.
type Manager interface {
	// Remove deletes an item from the filesystem.
	Remove(ctx context.Context, path *URI) error
	// Copy duplicates an item from a source path to a destination path.
	Copy(ctx context.Context, src, dst *URI, opts CopyOptions) error
	// Move relocates an item from a source path to a destination path.
	Move(ctx context.Context, src, dst *URI) error
}
