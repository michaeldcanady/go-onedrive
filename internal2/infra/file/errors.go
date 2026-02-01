package file

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrNotFolder    = errors.New("not a folder")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
	ErrPrecondition = errors.New("precondition error")
	ErrTransient    = errors.New("transient")
)

type DomainError struct {
	Kind    error
	DriveID string
	Path    string
	Err     error
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %v", e.Kind, e.Err)
}

func (e *DomainError) Unwrap() error { return e.Err }
