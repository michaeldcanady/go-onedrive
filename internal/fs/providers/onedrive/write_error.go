package onedrive

import "fmt"

// WriteError represents an error that occurred while writing to OneDrive.
type WriteError struct {
	// Path is the path that was being written.
	Path string
	// Err is the underlying error.
	Err error
}

// NewWriteError creates a new WriteError.
func NewWriteError(path string, err error) *WriteError {
	return &WriteError{
		Path: path,
		Err:  err,
	}
}

// Error returns a formatted error message.
func (e *WriteError) Error() string {
	return fmt.Sprintf("failed to write to path %s: %v", e.Path, e.Err)
}

// Unwrap returns the underlying error.
func (e *WriteError) Unwrap() error {
	return e.Err
}
