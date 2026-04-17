package mount

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the mount operation.
type Command struct {
	config config.Service
	log    logger.Logger
}

// NewCommand initializes a new instance of the mount Command.
func NewCommand(cfg config.Service, l logger.Logger) *Command {
	return &Command{
		config: cfg,
		log:    l,
	}
}

// Validate prepares and validates the options for the mount operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute adds a new mount point to the configuration.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("path", opts.Path),
		logger.String("type", opts.Type),
	)

	log.Info("adding mount point")

	m := config.MountConfig{
		Path:       opts.Path,
		Type:       opts.Type,
		IdentityID: opts.IdentityID,
		Options:    opts.MountOptions,
	}

	if err := c.config.AddMount(ctx, m); err != nil {
		log.Error("failed to add mount point", logger.Error(err))
		return fmt.Errorf("failed to add mount point %s: %w", opts.Path, err)
	}

	log.Info("mount point added successfully")
	fmt.Fprintf(opts.Stdout, "Mounted %s at %s\n", opts.Type, opts.Path)

	return nil
}

// Finalize performs any necessary cleanup after the mount operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
