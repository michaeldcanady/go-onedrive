package formatting

import "fmt"

// UnsupportedFormatError indicates that a requested output format is not supported.
type UnsupportedFormatError struct {
	// Format is the format identifier that is not supported.
	Format Format
	// Err is the underlying error, if any.
	Err error
}

func NewUnsupportedFormatError(format Format, err error) *UnsupportedFormatError {
	return &UnsupportedFormatError{
		Format: format,
		Err:    err,
	}
}

func (e *UnsupportedFormatError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("unsupported output format %s: %v", e.Format, e.Err)
	}
	return fmt.Sprintf("unsupported output format: %s", e.Format)
}

func (e *UnsupportedFormatError) Unwrap() error {
	return e.Err
}

// TableConfigurationError indicates that a table formatter was incorrectly configured.
type TableConfigurationError struct {
	// Reason is the human-readable explanation of the configuration issue.
	Reason string
	// Err is the underlying error, if any.
	Err error
}

func NewTableConfigurationError(reason string, err error) *TableConfigurationError {
	return &TableConfigurationError{
		Reason: reason,
		Err:    err,
	}
}

func (e *TableConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("invalid table configuration: %s: %v", e.Reason, e.Err)
	}
	return fmt.Sprintf("invalid table configuration: %s", e.Reason)
}

func (e *TableConfigurationError) Unwrap() error {
	return e.Err
}
