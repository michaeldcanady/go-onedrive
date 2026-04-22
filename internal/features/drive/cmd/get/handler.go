package get

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Command executes the drive get operation.
type Command struct {
	drive drive.Service
	log   logger.Logger
}

// NewCommand initializes a new instance of the drive get Command.
func NewCommand(d drive.Service, l logger.Logger) *Command {
	return &Command{
		drive: d,
		log:   l,
	}
}

// Validate prepares and validates the options for the drive get operation.
func (c *Command) Validate(ctx *CommandContext) error {
	return ctx.Options.Validate()
}

// Execute retrieves and displays the personal drive.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Debug("fetching personal drive", logger.String("identity", ctx.Options.IdentityID))
	d, err := c.drive.ResolvePersonalDrive(ctx.Ctx, ctx.Options.IdentityID)
	if err != nil {
		log.Error("failed to get personal drive", logger.Error(err))
		return fmt.Errorf("failed to get personal drive: %w", err)
	}

	log.Info("personal drive retrieved successfully", logger.String("id", d.ID))
	fmt.Fprintf(ctx.Options.Stdout, "Personal drive: %s (%s)\n", d.Name, d.ID)

	return nil
}

// Finalize performs any necessary cleanup after the drive get operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
