package remove

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive alias remove operation.
type Command struct {
	alias alias.Service
	log   logger.Logger
}

// NewCommand initializes a new instance of the drive alias remove Command.
func NewCommand(a alias.Service, l logger.Logger) *Command {
	return &Command{
		alias: a,
		log:   l,
	}
}

// Validate prepares and validates the options for the drive alias remove operation.
func (c *Command) Validate(ctx *CommandContext) error {
	return nil
}

// Execute deletes a drive alias.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("alias", ctx.Options.Alias),
	)

	log.Info("removing drive alias")

	if err := c.alias.DeleteAlias(ctx.Ctx, ctx.Options.Alias); err != nil {
		log.Error("failed to remove alias", logger.Error(err))
		return fmt.Errorf("failed to remove alias %s: %w", ctx.Options.Alias, err)
	}

	log.Info("alias removed successfully")
	fmt.Fprintf(ctx.Options.Stdout, "Alias %s removed successfully.\n", ctx.Options.Alias)

	return nil
}

// Finalize performs any necessary cleanup after the drive alias remove operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
