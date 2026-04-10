package onedrive

import "fmt"

// BadPathError represents an error where a path is malformed for OneDrive.
type BadPathError struct {
	// Path is the malformed path.
	Path string
	// Reason is the human-readable explanation of why the path is malformed.
	Reason string
	// Err is the underlying error, if any.
	Err error
}

// NewBadPathError creates a new BadPathError.
func NewBadPathError(path, reason string, err error) *BadPathError {
	return &BadPathError{
		Path:   path,
		Reason: reason,
		Err:    err,
	}
}

// Error returns a formatted error message.
func (e *BadPathError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("path %s is malformed: %s: %v", e.Path, e.Reason, e.Err)
	}
	return fmt.Sprintf("path %s is malformed: %s", e.Path, e.Reason)
}

// Unwrap returns the underlying error.
func (e *BadPathError) Unwrap() error {
	return e.Err
}
