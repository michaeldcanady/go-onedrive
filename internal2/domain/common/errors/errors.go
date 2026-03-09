package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound indicates the item (file or folder) was not found.
	ErrNotFound = errors.New("not found")
	// ErrNotFolder indicates the specified item is not a folder.
	ErrNotFolder = errors.New("not a folder")
	// ErrUnauthorized indicates the user is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates the user is authenticated but does not have
	// permission for the operation.
	ErrForbidden = errors.New("forbidden")
	// ErrConflict indicates the operation failed because of a resource conflict,
	// such as a file with the same name.
	ErrConflict = errors.New("conflict")
	// ErrInternal indicates an unexpected internal error.
	ErrInternal = errors.New("internal error")
	// ErrPrecondition indicates a precondition (like ETag) check failed.
	ErrPrecondition = errors.New("precondition error")
	// ErrTransient indicates a temporary error that can be retried.
	ErrTransient = errors.New("transient")
	// ErrInvalidRequest indicates the request was malformed.
	ErrInvalidRequest = errors.New("invalid request")
)

type DomainError struct {
	Kind    error
	Err     error
	DriveID string
	Path    string
}

func (e *DomainError) Error() string {
	msg := e.Kind.Error()
	if e.Err != nil {
		msg = fmt.Sprintf("%v: %v", e.Kind, e.Err)
	}
	if e.DriveID != "" || e.Path != "" {
		msg = fmt.Sprintf("%s (DriveID: %s, Path: %s)", msg, e.DriveID, e.Path)
	}
	return msg
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func (e *DomainError) Is(target error) bool {
	return e.Kind == target || errors.Is(e.Err, target)
}
