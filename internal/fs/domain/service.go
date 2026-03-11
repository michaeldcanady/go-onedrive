// Package domain defines the core domain models and interfaces for the OneDrive filesystem.
package domain

import (
	"context"
	"io"
)

// Reader defines the methods for reading data from the OneDrive filesystem.
type Reader interface {
	// Get retrieves a single item by its path.
	Get(ctx context.Context, path string) (Item, error)
	// List returns the children of a directory at the specified path.
	List(ctx context.Context, path string, opts ListOptions) ([]Item, error)
	// Stat returns metadata about an item at the specified path.
	Stat(ctx context.Context, path string, opts StatOptions) (Item, error)
	// ReadFile returns an io.ReadCloser for the contents of a file.
	ReadFile(ctx context.Context, path string, opts ReadOptions) (io.ReadCloser, error)
}

// Writer defines the methods for creating and updating items in the OneDrive filesystem.
type Writer interface {
	// WriteFile writes data to a file in the OneDrive filesystem.
	WriteFile(ctx context.Context, path string, r io.Reader, opts WriteOptions) (Item, error)
	// Mkdir creates a new directory at the specified path.
	Mkdir(ctx context.Context, path string, opts MKDirOptions) error
	// Upload uploads a file from a source path (local or remote) to the specified destination.
	Upload(ctx context.Context, src, dst string, opts UploadOptions) (Item, error)
	// Touch creates a new empty file or updates the modified time of an existing one.
	Touch(ctx context.Context, path string, opts TouchOptions) (Item, error)
}

// Manager defines the methods for managing items in the OneDrive filesystem.
type Manager interface {
	// Remove deletes an item at the specified path.
	Remove(ctx context.Context, path string, opts RemoveOptions) error
	// Copy copies an item from one path to another.
	Copy(ctx context.Context, src, dst string, opts CopyOptions) error
	// Move moves or renames an item from one path to another.
	Move(ctx context.Context, src, dst string, opts MoveOptions) error
}

// Service provides a full-featured interface for interacting with the OneDrive filesystem.
// It is a composition of Reader, Writer, and Manager interfaces.
type Service interface {
	Reader
	Writer
	Manager
}
