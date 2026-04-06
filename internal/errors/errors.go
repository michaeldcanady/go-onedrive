package errors

import (
	"errors"
	"fmt"
)

// Base error kinds
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
	ErrTransient = errors.New("transient error")
	// ErrInvalidRequest indicates the request was malformed.
	ErrInvalidRequest = errors.New("invalid request")
)

// DomainError represents a domain-level error with additional context.
type DomainError struct {
	Kind    error
	Err     error
	DriveID string
	Path    string
}

// Error returns the error message string.
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

// Unwrap returns the underlying error.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is reports whether the error matches the target.
func (e *DomainError) Is(target error) bool {
	return e.Kind == target || errors.Is(e.Err, target)
}

// NewDomainError creates a new DomainError.
func NewDomainError(kind error, err error, path string) *DomainError {
	return &DomainError{
		Kind: kind,
		Err:  err,
		Path: path,
	}
}

// CLIError represents a UI-level error intended for the end user.
type CLIError struct {
	Message  string
	ExitCode int
	Wrapped  error
}

// Error returns the user-facing error message.
func (e *CLIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}
	return "an unknown error occurred"
}

// Unwrap returns the underlying error.
func (e *CLIError) Unwrap() error {
	return e.Wrapped
}

// Exit codes for common scenarios
const (
	ExitSuccess       = 0
	ExitError         = 1
	ExitInvalidUsage  = 2
	ExitNotFound      = 4
	ExitConflict      = 5
	ExitUnauthorized  = 6
	ExitForbidden     = 7
	ExitTransient     = 8
	ExitInternalError = 9
)

// MapToCLI converts a standard error or DomainError into a CLIError.
func MapToCLI(err error) *CLIError {
	if err == nil {
		return nil
	}

	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr
	}

	var de *DomainError
	if errors.As(err, &de) {
		return mapDomainToCLI(de)
	}

	// Check against base error kinds directly
	code := ExitError
	switch {
	case errors.Is(err, ErrNotFound):
		code = ExitNotFound
	case errors.Is(err, ErrConflict):
		code = ExitConflict
	case errors.Is(err, ErrUnauthorized):
		code = ExitUnauthorized
	case errors.Is(err, ErrForbidden):
		code = ExitForbidden
	case errors.Is(err, ErrTransient):
		code = ExitTransient
	case errors.Is(err, ErrInternal):
		code = ExitInternalError
	case errors.Is(err, ErrInvalidRequest):
		code = ExitInvalidUsage
	}

	return &CLIError{
		Message:  err.Error(),
		ExitCode: code,
		Wrapped:  err,
	}
}

func mapDomainToCLI(de *DomainError) *CLIError {
	msg := de.Error()
	code := ExitError

	switch {
	case errors.Is(de.Kind, ErrNotFound):
		code = ExitNotFound
	case errors.Is(de.Kind, ErrConflict):
		code = ExitConflict
	case errors.Is(de.Kind, ErrUnauthorized):
		code = ExitUnauthorized
	case errors.Is(de.Kind, ErrForbidden):
		code = ExitForbidden
	case errors.Is(de.Kind, ErrTransient):
		code = ExitTransient
	case errors.Is(de.Kind, ErrInternal):
		code = ExitInternalError
	case errors.Is(de.Kind, ErrInvalidRequest):
		code = ExitInvalidUsage
	}

	return &CLIError{
		Message:  msg,
		ExitCode: code,
		Wrapped:  de,
	}
}
