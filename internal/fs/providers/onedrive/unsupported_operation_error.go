package onedrive

import "fmt"

// UnsupportedOperationError represents an error where an operation is not supported by the OneDrive provider.
type UnsupportedOperationError struct {
	// Path is the path that was involved in the operation.
	Path string
	// Operation is the name of the operation that is not supported.
	Operation string
	// Err is the underlying error, if any.
	Err error
}

// NewUnsupportedOperationError creates a new UnsupportedOperationError.
func NewUnsupportedOperationError(path, operation string, err error) *UnsupportedOperationError {
	return &UnsupportedOperationError{
		Path:      path,
		Operation: operation,
		Err:       err,
	}
}

// Error returns a formatted error message.
func (e *UnsupportedOperationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("operation %s is not supported for path %s: %v", e.Operation, e.Path, e.Err)
	}
	return fmt.Sprintf("operation %s is not supported for path %s", e.Operation, e.Path)
}

// Unwrap returns the underlying error.
func (e *UnsupportedOperationError) Unwrap() error {
	return e.Err
}
