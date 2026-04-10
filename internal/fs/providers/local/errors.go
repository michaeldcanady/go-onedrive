package local

import "fmt"

// ReadError represents an error that occurred while reading from the local filesystem.
type ReadError struct {
	Path string
	Err  error
}

func NewReadError(path string, err error) *ReadError {
	return &ReadError{
		Path: path,
		Err:  err,
	}
}

func (e *ReadError) Error() string {
	return fmt.Sprintf("failed to read from local path %s: %v", e.Path, e.Err)
}

func (e *ReadError) Unwrap() error {
	return e.Err
}

// WriteError represents an error that occurred while writing to the local filesystem.
type WriteError struct {
	Path string
	Err  error
}

func NewWriteError(path string, err error) *WriteError {
	return &WriteError{
		Path: path,
		Err:  err,
	}
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("failed to write to local path %s: %v", e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}
