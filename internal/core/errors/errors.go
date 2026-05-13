// Package errors provides common error variables and utilities for error handling throughout the application.
package errors

import (
	"errors"
	"fmt"
)

// Common error variables used for consistent error reporting across domain boundaries.
var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInternal         = errors.New("internal error")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidPath      = errors.New("invalid path")
	ErrNotEmpty         = errors.New("not empty")
	ErrUnavailable      = errors.New("unavailable")
)

// Wrap returns an error with the provided message prepended to the original error's message.
// It returns nil if the provided error is nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Is reports whether any error in the chain matches the target.
// It is a wrapper around the standard library's [errors.Is].
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in the chain that matches the target and sets target to that error value.
// It is a wrapper around the standard library's [errors.As].
func As(err error, target any) bool {
	return errors.As(err, target)
}
