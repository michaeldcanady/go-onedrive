package util

import (
	"errors"
	"fmt"
	"io"

	"github.com/fatih/color"
	domainerrors "github.com/michaeldcanady/go-onedrive/internal2/domain/common/errors"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
)

// BaseCommand provides common functionality for all CLI commands.
type BaseCommand struct {
	Container di.Container
	Log       logger.Logger
	Name      string
}

// NewBaseCommand creates a new BaseCommand.
func NewBaseCommand(container di.Container, name string) BaseCommand {
	return BaseCommand{
		Container: container,
		Name:      name,
	}
}

// WithLogger allows injecting a logger into BaseCommand.
func (c *BaseCommand) WithLogger(log logger.Logger) *BaseCommand {
	c.Log = log
	return c
}

// Initialize ensures the logger is set up for the command.
func (c *BaseCommand) Initialize(id string) error {
	if c.Log != nil {
		return nil
	}

	l, err := EnsureLogger(c.Container, id)
	if err != nil {
		return NewCommandError(c.Name, "failed to initialize logger", err)
	}
	c.Log = l
	return nil
}

// RenderError renders a domain error in a user-friendly way.
func (c *BaseCommand) RenderError(w io.Writer, err error) {
	if err == nil {
		return
	}

	red := color.New(color.FgRed, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	var domainErr *domainerrors.DomainError
	if errors.As(err, &domainErr) {
		fmt.Fprintf(w, "%s %s\n", red("Error:"), domainErr.Error())
		if domainErr.Err != nil {
			fmt.Fprintf(w, "  %s %v\n", yellow("Details:"), domainErr.Err)
		}
		return
	}

	// Fallback for regular errors
	fmt.Fprintf(w, "%s %v\n", red("Error:"), err)
}
