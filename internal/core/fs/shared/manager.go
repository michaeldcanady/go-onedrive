package shared

import (
	"context"
)

// Manager defines the operations for managing the structure and lifecycle of filesystem items.
type Manager interface {
	// Remove deletes the item at the specified path.
	Remove(ctx context.Context, path string) error
	// Copy duplicates an item from a source path to a destination path.
	Copy(ctx context.Context, src, dst string, opts CopyOptions) error
	// Move relocates or renames an item within the filesystem.
	Move(ctx context.Context, src, dst string) error
}
