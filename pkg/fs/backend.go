package fs

import (
	"context"
	"io"
)

// Backend defines the primitive operations for a storage engine.
// This interface is designed to be easily mapped to a gRPC service contract.
type Backend interface {
	Namer

	// Stat returns metadata for an item at the specified path.
	Stat(ctx context.Context, path string) (Item, error)

	// List returns the immediate children of the specified directory path.
	List(ctx context.Context, path string) ([]Item, error)

	// Open returns an io.ReadCloser for the content of the file at the specified path.
	Open(ctx context.Context, path string) (io.ReadCloser, error)

	// Create uploads or updates a file with the content from the provided reader at the specified path.
	Create(ctx context.Context, path string, r io.Reader) (Item, error)

	// Mkdir creates a new directory at the specified path.
	Mkdir(ctx context.Context, path string) error

	// Remove deletes an item from the storage at the specified path.
	Remove(ctx context.Context, path string) error

	// Capabilities returns the advanced operations supported by this backend.
	Capabilities() Capabilities
}

// Capabilities defines the set of advanced operations a backend might support natively.
type Capabilities struct {
	// CanMove indicates if the backend can perform native renames/moves.
	CanMove bool
	// CanCopy indicates if the backend can perform native server-side copies.
	CanCopy bool
	// CanRecursive indicates if the backend can perform recursive operations natively.
	CanRecursive bool
}

// AdvancedBackend is an optional interface for backends that support native high-level operations.
type AdvancedBackend interface {
	Backend
	// Move performs a native move/rename.
	Move(ctx context.Context, src, dst string) error
	// Copy performs a native copy.
	Copy(ctx context.Context, src, dst string) error
}
