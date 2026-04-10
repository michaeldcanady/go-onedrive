package filtering

import "fmt"

// UnknownItemTypeError indicates that a filter was requested for an invalid or unknown item type.
type UnknownItemTypeError struct {
	// Err is the underlying error, if any.
	Err error
}

func NewUnknownItemTypeError(err error) *UnknownItemTypeError {
	return &UnknownItemTypeError{
		Err: err,
	}
}

func (e *UnknownItemTypeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("filtered item type is unknown: %v", e.Err)
	}
	return "filtered item type is unknown"
}

func (e *UnknownItemTypeError) Unwrap() error {
	return e.Err
}
