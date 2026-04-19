package use

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/shared"
)

// Command executes the profile use operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile use Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile use operation.
func (c *Command) Validate(ctx *CommandContext) error {
	ctx.Options.Name = strings.TrimSpace(ctx.Options.Name)
	if ctx.Options.Name == "" {
		c.log.Error("profile name is missing")
		return errors.New("profile name is required")
	}
	if ctx.Options.Stdout == nil {
		c.log.Error("stdout is missing")
		return errors.New("stdout is required")
	}
	return nil
}

// Execute sets the active profile.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("switching active profile", logger.String("name", ctx.Options.Name))
	if err := c.profile.SetActive(ctx.Ctx, ctx.Options.Name, shared.ScopeGlobal); err != nil {
		log.Error("failed to switch profile", logger.String("name", ctx.Options.Name), logger.Error(err))
		return fmt.Errorf("failed to switch to profile %s: %w", ctx.Options.Name, err)
	}

	log.Info("switched active profile successfully", logger.String("name", ctx.Options.Name))
	return nil
}

// Finalize performs any necessary cleanup after the profile use operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	_, _ = fmt.Fprintf(ctx.Options.Stdout, "Switched to profile '%s'.\n", ctx.Options.Name)
	return nil
}
