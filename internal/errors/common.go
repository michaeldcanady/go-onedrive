package errors

import "fmt"

// NotFoundError represents a resource that could not be found.
type NotFoundError struct {
	Path string
	Err  error
}

func NewNotFoundError(path string, err error) *NotFoundError {
	return &NotFoundError{Path: path, Err: err}
}

func (e *NotFoundError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("resource not found at path %s", e.Path)
	}
	return "resource not found"
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// ForbiddenError represents an operation that was denied due to lack of permissions.
type ForbiddenError struct {
	Path string
	Err  error
}

func NewForbiddenError(path string, err error) *ForbiddenError {
	return &ForbiddenError{Path: path, Err: err}
}

func (e *ForbiddenError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("access denied for resource at %s", e.Path)
	}
	return "access denied"
}

func (e *ForbiddenError) Unwrap() error {
	return e.Err
}

// UnauthorizedError represents a failure due to missing or invalid authentication.
type UnauthorizedError struct {
	Err error
}

func NewUnauthorizedError(err error) *UnauthorizedError {
	return &UnauthorizedError{Err: err}
}

func (e *UnauthorizedError) Error() string {
	return "unauthorized: authentication is missing or invalid"
}

func (e *UnauthorizedError) Unwrap() error {
	return e.Err
}

// ConflictError represents a conflict with the current state of the resource (e.g., already exists).
type ConflictError struct {
	Path string
	Err  error
}

func NewConflictError(path string, err error) *ConflictError {
	return &ConflictError{Path: path, Err: err}
}

func (e *ConflictError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("resource conflict at path %s", e.Path)
	}
	return "resource conflict"
}

func (e *ConflictError) Unwrap() error {
	return e.Err
}

// BadRequestError represents a malformed request or invalid parameters.
type BadRequestError struct {
	Err error
}

func NewBadRequestError(err error) *BadRequestError {
	return &BadRequestError{Err: err}
}

func (e *BadRequestError) Error() string {
	return "bad request: the request was malformed or contained invalid parameters"
}

func (e *BadRequestError) Unwrap() error {
	return e.Err
}

// PreconditionFailedError represents a failure of a precondition (e.g., ETag mismatch).
type PreconditionFailedError struct {
	Path string
	Err  error
}

func NewPreconditionFailedError(path string, err error) *PreconditionFailedError {
	return &PreconditionFailedError{Path: path, Err: err}
}

func (e *PreconditionFailedError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("precondition failed for path %s", e.Path)
	}
	return "precondition failed"
}

func (e *PreconditionFailedError) Unwrap() error {
	return e.Err
}

// TransientError represents temporary failures that can be retried.
type TransientError struct {
	Reason string
	Err    error
}

func NewTransientError(reason string, err error) *TransientError {
	return &TransientError{Reason: reason, Err: err}
}

func (e *TransientError) Error() string {
	return fmt.Sprintf("transient error: %s", e.Reason)
}

func (e *TransientError) Unwrap() error {
	return e.Err
}

// InternalError represents an unexpected internal failure.
type InternalError struct {
	Err error
}

func NewInternalError(err error) *InternalError {
	return &InternalError{Err: err}
}

func (e *InternalError) Error() string {
	return "internal error: an unexpected failure occurred"
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

// IllegalCharacterError represents an error where a path contains characters that are not allowed.
type IllegalCharacterError struct {
	// Path is the path that contains the illegal character.
	Path string
	// Character is the illegal character found.
	Character string
	// Err is the underlying error, if any.
	Err error
}

// NewIllegalCharacterError creates a new IllegalCharacterError.
func NewIllegalCharacterError(path, character string, err error) *IllegalCharacterError {
	return &IllegalCharacterError{
		Path:      path,
		Character: character,
		Err:       err,
	}
}

// Error returns a formatted error message.
func (e *IllegalCharacterError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("path %s contains illegal character %s: %v", e.Path, e.Character, e.Err)
	}
	return fmt.Sprintf("path %s contains illegal character %s", e.Path, e.Character)
}

// Unwrap returns the underlying error.
func (e *IllegalCharacterError) Unwrap() error {
	return e.Err
}

// TrailingSlashError represents an error where a path contains a trailing slash, which is not allowed for files.
type TrailingSlashError struct {
	// Path is the path that contains the trailing slash.
	Path string
	// Err is the underlying error, if any.
	Err error
}

// NewTrailingSlashError creates a new TrailingSlashError.
func NewTrailingSlashError(path string, err error) *TrailingSlashError {
	return &TrailingSlashError{
		Path: path,
		Err:  err,
	}
}

// Error returns the error message for the TrailingSlashError.
func (e *TrailingSlashError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("path %s contains a trailing slash: %v", e.Path, e.Err)
	}
	return fmt.Sprintf("path %s contains a trailing slash", e.Path)
}

// Unwrap returns the underlying error.
func (e *TrailingSlashError) Unwrap() error {
	return e.Err
}

// UnregisteredProviderError indicates that a provider is not registered in the registry.
type UnregisteredProviderError struct {
	// Provider is the name of the provider that is not registered.
	Provider string
	// Err is the underlying error, if any.
	Err error
}

// NewUnregisteredProviderError creates a new UnregisteredProviderError.
func NewUnregisteredProviderError(provider string, err error) *UnregisteredProviderError {
	return &UnregisteredProviderError{
		Provider: provider,
		Err:      err,
	}
}

// Error returns a formatted error message.
func (e *UnregisteredProviderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("provider '%s' is not registered: %v", e.Provider, e.Err)
	}
	return fmt.Sprintf("provider '%s' is not registered", e.Provider)
}

// Unwrap returns the underlying error.
func (e *UnregisteredProviderError) Unwrap() error {
	return e.Err
}
