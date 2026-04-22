package delete

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Command executes the profile delete operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile delete Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile delete operation.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Options.Stdout == nil {
		return fmt.Errorf("stdout must not be nil")
	}

	ctx.Options.Name = strings.TrimSpace(ctx.Options.Name)
	if ctx.Options.Name == "" {
		return errors.New("profile name is required")
	}

	return nil
}

// Execute deletes a profile.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("deleting profile", logger.String("name", ctx.Options.Name))
	if err := c.profile.Delete(ctx.Ctx, ctx.Options.Name); err != nil {
		log.Error("failed to delete profile", logger.String("name", ctx.Options.Name), logger.Error(err))
		return fmt.Errorf("failed to delete profile %s: %w", ctx.Options.Name, err)
	}

	log.Info("profile deleted successfully", logger.String("name", ctx.Options.Name))
	return nil
}

// Finalize performs any necessary cleanup after the profile delete operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	_, _ = fmt.Fprintf(ctx.Options.Stdout, "Profile '%s' deleted successfully.\n", ctx.Options.Name)
	return nil
}
