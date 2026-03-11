package login

import (
	"strings"
)

type errorUnwrapper interface {
	Unwrap() error
}

var _ (errorUnwrapper) = (*CommandError)(nil)
var _ (error) = (*CommandError)(nil)

type CommandError struct {
	CommandName string
	Message     string
	NestedError error
}

func NewCommandErrorWithName(name string) *CommandError {
	return NewCommandErrorWithNameWithMessage(name, "")
}

func NewCommandErrorWithNameWithError(name string, err error) *CommandError {
	return NewCommandError(name, "", err)
}

func NewCommandErrorWithNameWithMessage(name, message string) *CommandError {
	return NewCommandError(name, message, nil)
}

func NewCommandError(name, message string, err error) *CommandError {
	return &CommandError{
		CommandName: name,
		Message:     message,
		NestedError: err,
	}
}

// Error implements [error].
func (c *CommandError) Error() string {
	builder := &strings.Builder{}
	builder.Write([]byte(c.CommandName))
	if strings.TrimSpace(c.Message) == "" {
		builder.Write([]byte(" "))
		builder.Write([]byte(c.Message))
	}
	if c.NestedError != nil {
		builder.Write([]byte(" "))
		builder.Write([]byte(c.NestedError.Error()))
	}
	return builder.String()
}

// Unwrap implements [errorUnwrapper].
func (c *CommandError) Unwrap() error {
	return c.NestedError
}
