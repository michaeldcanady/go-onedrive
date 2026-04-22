package set

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Command executes the config set operation.
type Command struct {
	config config.Service
	log    logger.Logger
}

// NewCommand initializes a new instance of the config set Command.
func NewCommand(c config.Service, l logger.Logger) *Command {
	return &Command{
		config: c,
		log:    l,
	}
}

// Validate prepares and validates the options for the config set operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute updates a configuration setting.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Debug("updating configuration setting", logger.String("key", opts.Key), logger.String("value", opts.Value))
	if err := c.config.UpdateConfig(ctx, opts.Key, opts.Value); err != nil {
		log.Error("failed to update configuration", logger.Error(err))
		return fmt.Errorf("failed to update configuration: %w", err)
	}

	log.Info("configuration updated successfully", logger.String("key", opts.Key))
	fmt.Fprintf(opts.Stdout, "Set %s to %s successfully.\n", opts.Key, opts.Value)

	return nil
}

// Finalize performs any necessary cleanup after the config set operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
