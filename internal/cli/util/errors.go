// Package util provides common utilities and base structures for the odc CLI.
package util

import (
	"errors"
	"strings"
)

// ErrUserAbort is returned when a user aborts an interactive prompt or sends a SIGINT (Ctrl+C).
var ErrUserAbort = errors.New("user aborted")

// errorUnwrapper is an interface for errors that can be unwrapped to reveal a nested error.
type errorUnwrapper interface {
	Unwrap() error
}

var _ (errorUnwrapper) = (*CommandError)(nil)
var _ (error) = (*CommandError)(nil)

// CommandError represents an error that occurs during the execution of a CLI command.
// It includes the command name, an optional user-facing message, and an optional nested error.
type CommandError struct {
	// CommandName is the name of the command that encountered the error.
	CommandName string
	// Message is a user-friendly description of the error.
	Message string
	// NestedError is the underlying error that caused the command to fail.
	NestedError error
}

// NewCommandErrorWithName creates a new CommandError with only the command name.
func NewCommandErrorWithName(name string) *CommandError {
	return NewCommandErrorWithNameWithMessage(name, "")
}

// NewCommandErrorWithNameWithError creates a new CommandError with a command name and an underlying error.
func NewCommandErrorWithNameWithError(name string, err error) *CommandError {
	return NewCommandError(name, "", err)
}

// NewCommandErrorWithNameWithMessage creates a new CommandError with a command name and a custom message.
func NewCommandErrorWithNameWithMessage(name, message string) *CommandError {
	return NewCommandError(name, message, nil)
}

// NewCommandError creates a fully initialized CommandError.
func NewCommandError(name, message string, err error) *CommandError {
	return &CommandError{
		CommandName: name,
		Message:     message,
		NestedError: err,
	}
}

// Error implements the [error] interface, providing a formatted string representation
// of the command name, message, and nested error.
func (c *CommandError) Error() string {
	builder := &strings.Builder{}
	builder.Write([]byte(c.CommandName))
	if strings.TrimSpace(c.Message) != "" {
		builder.Write([]byte(": "))
		builder.Write([]byte(c.Message))
	}
	if c.NestedError != nil {
		builder.Write([]byte(": "))
		builder.Write([]byte(c.NestedError.Error()))
	}
	return builder.String()
}

// Unwrap implements the [errorUnwrapper] interface, allowing access to the
// underlying error.
func (c *CommandError) Unwrap() error {
	return c.NestedError
}
