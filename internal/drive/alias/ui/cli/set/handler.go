package set

import (
	"context"
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
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return nil
}

// Execute defines or updates a drive alias.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("alias", opts.Alias),
		logger.String("id", opts.DriveID),
	)

	log.Info("setting drive alias")

	if err := c.alias.SetAlias(opts.DriveID, opts.Alias); err != nil {
		log.Error("failed to set alias", logger.Error(err))
		return fmt.Errorf("failed to set alias %s -> %s: %w", opts.Alias, opts.DriveID, err)
	}

	log.Info("alias set successfully")
	fmt.Fprintf(opts.Stdout, "Alias %s set to %s successfully.\n", opts.Alias, opts.DriveID)

	return nil
}

// Finalize performs any necessary cleanup after the drive alias set operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
