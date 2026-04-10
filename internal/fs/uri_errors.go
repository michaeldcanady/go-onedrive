package fs

import "fmt"

// InvalidURIError indicates that a string could not be parsed as a valid URI.
type InvalidURIError struct {
	// URI is the raw string that failed to parse.
	URI string
	// Reason is the human-readable explanation of why the URI is invalid.
	Reason string
	// Err is the underlying error, if any.
	Err error
}

func NewInvalidURIError(uri, reason string, err error) *InvalidURIError {
	return &InvalidURIError{
		URI:    uri,
		Reason: reason,
		Err:    err,
	}
}

func (e *InvalidURIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("invalid URI '%s': %s: %v", e.URI, e.Reason, e.Err)
	}
	return fmt.Sprintf("invalid URI '%s': %s", e.URI, e.Reason)
}

func (e *InvalidURIError) Unwrap() error {
	return e.Err
}
