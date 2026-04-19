package editor

import (
	"context"
	"io"
)

// Service defines the interface for editor-related operations.
type Service interface {
	// Launch launches an external editor with the specified path.
	Launch(ctx context.Context, path string) error
	// LaunchTempFile creates a temporary file with the specified prefix and suffix,
	// writes the content from the reader to it, launches the editor,
	// and returns the modified content.
	LaunchTempFile(ctx context.Context, prefix, suffix string, reader io.Reader) ([]byte, error)
	// WithOptions returns a new Service instance with the specified options applied.
	WithOptions(opts ...Option) Service
}
