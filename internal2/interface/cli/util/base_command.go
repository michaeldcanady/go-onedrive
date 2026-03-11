package util

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	domainerrors "github.com/michaeldcanady/go-onedrive/internal2/domain/common/errors"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
)

// BaseCommand provides common functionality for all CLI commands.
type BaseCommand struct {
	Container di.Container
	Log       logger.Logger
	Name      string
	Quiet     bool
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

// WithQuiet allows setting the quiet flag.
func (c *BaseCommand) WithQuiet(quiet bool) *BaseCommand {
	c.Quiet = quiet
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

// RenderWarning renders a warning message in a user-friendly way.
func (c *BaseCommand) RenderWarning(w io.Writer, format string, a ...any) {
	if c.Quiet {
		return
	}
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(w, "%s %s\n", yellow("Warning:"), msg)
}

// RenderInfo renders an informational message in a user-friendly way.
func (c *BaseCommand) RenderInfo(w io.Writer, format string, a ...any) {
	if c.Quiet {
		return
	}
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(w, "%s %s\n", blue("Info:"), msg)
}

// RenderSuccess renders a success message in a user-friendly way.
func (c *BaseCommand) RenderSuccess(w io.Writer, format string, a ...any) {
	if c.Quiet {
		return
	}
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(w, "%s %s\n", green("Success:"), msg)
}

func (c *BaseCommand) Prompt(prompt promptui.Prompt) (string, error) {
	result, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrAbort) {
			// TODO: have it return custom abort error?
			return "", nil
		}
	}
	return result, err
}

// PromptConfirm asks the user for confirmation.
func (c *BaseCommand) PromptConfirm(w io.Writer, label string) (bool, error) {
	result, err := c.Prompt(promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Stdout:    NewNopWriteCloser(w),
	})
	if err != nil {
		return false, err
	}
	return parseBool(result)
}

func parseBool(str string) (bool, error) {
	switch str {
	case "yes", "y":
		return true, nil
	case "no", "n":
		return false, nil
	default:
		return strconv.ParseBool(str)
	}
}
