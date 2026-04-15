package use

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
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
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute sets the active OneDrive drive.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("drive_ref", opts.DriveRef),
	)

	log.Info("switching active drive")

	log.Debug("checking if drive reference is an alias")
	id, err := c.alias.GetDriveIDByAlias(opts.DriveRef)
	if err != nil {
		log.Debug("drive reference is not an alias, using as ID")
		id = opts.DriveRef
	}

	log.Debug("setting active drive", logger.String("id", id))
	if err := c.drive.SetActive(ctx, id, state.ScopeGlobal); err != nil {
		log.Error("failed to switch drive", logger.Error(err))
		return fmt.Errorf("failed to switch to drive %s: %w", opts.DriveRef, err)
	}

	log.Info("active drive switched successfully")
	fmt.Fprintf(opts.Stdout, "Switched to drive: %s\n", opts.DriveRef)

	return nil
}

// Finalize performs any necessary cleanup after the drive use operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
