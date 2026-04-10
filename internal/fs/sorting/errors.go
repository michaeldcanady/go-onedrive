package sorting

import "fmt"

// ComparisonError indicates that two values could not be compared.
type ComparisonError struct {
	// Reason is the human-readable explanation of why comparison failed.
	Reason string
	// Err is the underlying error, if any.
	Err error
}

func NewComparisonError(reason string, err error) *ComparisonError {
	return &ComparisonError{
		Reason: reason,
		Err:    err,
	}
}

func (e *ComparisonError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("comparison failed: %s: %v", e.Reason, e.Err)
	}
	return fmt.Sprintf("comparison failed: %s", e.Reason)
}

func (e *ComparisonError) Unwrap() error {
	return e.Err
}

// UnknownSortFieldError indicates that a sort was requested for a field that does not exist.
type UnknownSortFieldError struct {
	// Field is the name of the field that was not found.
	Field string
	// Err is the underlying error, if any.
	Err error
}

func NewUnknownSortFieldError(field string, err error) *UnknownSortFieldError {
	return &UnknownSortFieldError{
		Field: field,
		Err:   err,
	}
}

func (e *UnknownSortFieldError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("unknown sort field: %s: %v", e.Field, e.Err)
	}
	return fmt.Sprintf("unknown sort field: %s", e.Field)
}

func (e *UnknownSortFieldError) Unwrap() error {
	return e.Err
}
