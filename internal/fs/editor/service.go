package editor

import (
	"context"
	"io"
)

// Service defines the interface for editor-related operations.
type Service interface {
	// Launch launches an external editor with the specified path.
	Launch(path string) error
	// LaunchTempFile creates a temporary file with the specified prefix and suffix,
	// writes the content from the reader to it, launches the editor,
	// and returns the modified content and the path to the temporary file.
	LaunchTempFile(ctx context.Context, prefix, suffix string, reader io.Reader) ([]byte, string, error)
	// WithIO returns a new Service instance with the specified standard input, output, and error writers.
	WithIO(stdin io.Reader, stdout, stderr io.Writer) Service
}
