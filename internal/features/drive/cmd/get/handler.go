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

// Execute retrieves and displays the specified drive.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Debug("fetching drive", logger.String("ref", ctx.Options.DriveRef), logger.String("identity", ctx.Options.IdentityID))
	d, err := c.drive.ResolveDrive(ctx.Ctx, ctx.Options.DriveRef, ctx.Options.IdentityID)
	if err != nil {
		log.Error("failed to resolve drive", logger.String("ref", ctx.Options.DriveRef), logger.Error(err))
		return fmt.Errorf("failed to resolve drive %s: %w", ctx.Options.DriveRef, err)
	}

	log.Info("drive resolved successfully", logger.String("id", d.ID))
	fmt.Fprintf(ctx.Options.Stdout, "Drive: %s (%s) [%s]\n", d.Name, d.ID, d.Type)

	return nil
}

// Finalize performs any necessary cleanup after the drive get operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
