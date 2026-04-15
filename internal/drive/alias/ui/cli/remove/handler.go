package remove

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
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
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return nil
}

// Execute deletes a drive alias.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("alias", opts.Alias),
	)

	log.Info("removing drive alias")

	if err := c.alias.DeleteAlias(opts.Alias); err != nil {
		log.Error("failed to remove alias", logger.Error(err))
		return fmt.Errorf("failed to remove alias %s: %w", opts.Alias, err)
	}

	log.Info("alias removed successfully")
	fmt.Fprintf(opts.Stdout, "Alias %s removed successfully.\n", opts.Alias)

	return nil
}

// Finalize performs any necessary cleanup after the drive alias remove operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
