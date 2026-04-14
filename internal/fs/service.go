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
	// Get retrieves a single Item by its structured URI.
	Get(ctx context.Context, uri *URI) (Item, error)
	// List returns the immediate children of the specified directory URI.
	List(ctx context.Context, uri *URI, opts ListOptions) ([]Item, error)
	// ReadFile provides an io.ReadCloser for the content of the file at the specified URI.
	ReadFile(ctx context.Context, uri *URI, opts ReadOptions) (io.ReadCloser, error)
	// Stat returns metadata for an item at the specified URI.
	Stat(ctx context.Context, uri *URI) (Item, error)
}

// Writer defines the operations for modifying and creating items in the filesystem.
type Writer interface {
	// WriteFile uploads or updates a file with the content from the provided reader at the specified URI.
	WriteFile(ctx context.Context, uri *URI, r io.Reader, opts WriteOptions) (Item, error)
	// Mkdir creates a new directory at the specified URI.
	Mkdir(ctx context.Context, uri *URI) error
	// Touch creates a new empty file or updates the modification time of an existing one at the specified URI.
	Touch(ctx context.Context, uri *URI) (Item, error)
}

// Manager defines the operations for higher-level filesystem management and item manipulation.
type Manager interface {
	// Remove deletes an item from the filesystem at the specified URI.
	Remove(ctx context.Context, uri *URI) error
	// Copy duplicates an item from a source URI to a destination URI.
	Copy(ctx context.Context, src, dst *URI, opts CopyOptions) error
	// Move relocates an item from a source URI to a destination URI.
	Move(ctx context.Context, src, dst *URI) error
}
