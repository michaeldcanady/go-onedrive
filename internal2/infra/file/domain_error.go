package file

import "fmt"

// DomainError represents a failure in the infrastructure layer that wraps
// a Graph API error into a domain-specific category.
type DomainError struct {
	// Kind is the classified type of error.
	Kind error
	// DriveID is the ID of the drive where the error occurred.
	DriveID string
	// Path is the path of the item that caused the error.
	Path string
	// Err is the original underlying error.
	Err error
}

// Error returns a formatted string representation of the DomainError.
func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %v", e.Kind, e.Err)
}

// Unwrap returns the underlying error.
func (e *DomainError) Unwrap() error { return e.Err }
