package fs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrForbidden     = errors.New("forbidden")
	ErrConflict      = errors.New("conflict")
	ErrInternal      = errors.New("internal error")
	ErrInvalidRequest = errors.New("invalid request")
)

type Error struct {
	Kind error
	Err  error
	Path string
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Kind.Error() + ": " + e.Err.Error()
	}
	return e.Kind.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}
