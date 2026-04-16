package set

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive alias set operation.
type Command struct {
	alias alias.Service
	log   logger.Logger
}

// NewCommand initializes a new instance of the drive alias set Command.
func NewCommand(a alias.Service, l logger.Logger) *Command {
	return &Command{
		alias: a,
		log:   l,
	}
}

// Validate prepares and validates the options for the drive alias set operation.
func (c *Command) Validate(ctx *CommandContext) error {
	return nil
}

// Execute defines or updates a drive alias.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("alias", ctx.Options.Alias),
		logger.String("id", ctx.Options.DriveID),
	)

	log.Info("setting drive alias")

	if err := c.alias.SetAlias(ctx.Options.DriveID, ctx.Options.Alias); err != nil {
		log.Error("failed to set alias", logger.Error(err))
		return fmt.Errorf("failed to set alias %s -> %s: %w", ctx.Options.Alias, ctx.Options.DriveID, err)
	}

	log.Info("alias set successfully")
	fmt.Fprintf(ctx.Options.Stdout, "Alias %s set to %s successfully.\n", ctx.Options.Alias, ctx.Options.DriveID)

	return nil
}

// Finalize performs any necessary cleanup after the drive alias set operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
