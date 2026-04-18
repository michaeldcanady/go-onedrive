package use

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Command executes the drive use operation.
type Command struct {
	drive drive.Service
	alias alias.Service
	log   logger.Logger
}

// NewCommand initializes a new instance of the drive use Command.
func NewCommand(d drive.Service, a alias.Service, l logger.Logger) *Command {
	return &Command{
		drive: d,
		alias: a,
		log:   l,
	}
}

// Validate prepares and validates the options for the drive use operation.
func (c *Command) Validate(ctx *CommandContext) error {
	return ctx.Options.Validate()
}

// Execute sets the active OneDrive drive.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("drive_ref", ctx.Options.DriveRef),
	)

	log.Info("switching active drive")

	log.Debug("checking if drive reference is an alias")
	id, err := c.alias.GetDriveIDByAlias(ctx.Ctx, ctx.Options.DriveRef)
	if err != nil {
		log.Debug("drive reference is not an alias, using as ID")
		id = ctx.Options.DriveRef
	}

	log.Debug("setting active drive", logger.String("id", id), logger.String("identity", ctx.Options.IdentityID))
	if err := c.drive.SetActive(ctx.Ctx, id, ctx.Options.IdentityID, state.ScopeGlobal); err != nil {
		log.Error("failed to switch drive", logger.Error(err))
		return fmt.Errorf("failed to switch to drive %s: %w", ctx.Options.DriveRef, err)
	}

	log.Info("active drive switched successfully")
	fmt.Fprintf(ctx.Options.Stdout, "Switched to drive: %s\n", ctx.Options.DriveRef)

	return nil
}

// Finalize performs any necessary cleanup after the drive use operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
