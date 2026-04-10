package onedrive

import "fmt"

// ReadError represents an error that occurred while reading from OneDrive.
type ReadError struct {
	// Path is the path that was being read.
	Path string
	// Err is the underlying error.
	Err error
}

// NewReadError creates a new ReadError.
func NewReadError(path string, err error) *ReadError {
	return &ReadError{
		Path: path,
		Err:  err,
	}
}

// Error returns a formatted error message.
func (e *ReadError) Error() string {
	return fmt.Sprintf("failed to read from path %s: %v", e.Path, e.Err)
}

// Unwrap returns the underlying error.
func (e *ReadError) Unwrap() error {
	return e.Err
}
