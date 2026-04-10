package onedrive

import "fmt"

// GoneError represents a 410 Gone error from OneDrive.
type GoneError struct {
	Path string
	Err  error
}

func NewGoneError(path string, err error) *GoneError {
	return &GoneError{Path: path, Err: err}
}

func (e *GoneError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("resource at path %s is gone (permanently deleted)", e.Path)
	}
	return "resource is gone (permanently deleted)"
}

func (e *GoneError) Unwrap() error {
	return e.Err
}

// LockedError represents a 423 Locked error from OneDrive.
type LockedError struct {
	Path string
	Err  error
}

func NewLockedError(path string, err error) *LockedError {
	return &LockedError{Path: path, Err: err}
}

func (e *LockedError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("resource at path %s is locked", e.Path)
	}
	return "resource is locked"
}

func (e *LockedError) Unwrap() error {
	return e.Err
}

// InsufficientStorageError represents a 507 Insufficient Storage error from OneDrive.
type InsufficientStorageError struct {
	Err error
}

func NewInsufficientStorageError(err error) *InsufficientStorageError {
	return &InsufficientStorageError{Err: err}
}

func (e *InsufficientStorageError) Error() string {
	return "insufficient storage: your OneDrive is full"
}

func (e *InsufficientStorageError) Unwrap() error {
	return e.Err
}
