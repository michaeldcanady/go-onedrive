package fs

import "fmt"

// TrailingSlashError represents an error where a path contains a trailing slash, which is not allowed in OneDrive item paths (allowed for directories)
type TrailingSlashError struct {
	// Path is the path that contains the trailing slash.
	Path string
}

// NewTrailingSlashError creates a new TrailingSlashError for the given path.
func NewTrailingSlashError(path string) *TrailingSlashError {
	return &TrailingSlashError{
		Path: path,
	}
}

// Error returns the error message for the TrailingSlashError.
func (e *TrailingSlashError) Error() string {
	return fmt.Sprintf("path %s contains a trailing slash", e.Path)
}
