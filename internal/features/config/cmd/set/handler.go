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

	log.Debug("fetching current configuration")
	cfg, err := c.config.GetConfig(ctx)
	if err != nil {
		log.Error("failed to retrieve configuration", logger.Error(err))
		return fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	log.Debug("updating configuration setting", logger.String("key", opts.Key), logger.String("value", opts.Value))
	switch opts.Key {
	case "auth.provider":
		cfg.Auth.Provider = opts.Value
	case "auth.method":
		cfg.Auth.Method = opts.Value
	case "logging.format":
		cfg.Logging.Format = opts.Value
	default:
		return fmt.Errorf("configuration key not supported via CLI: %s", opts.Key)
	}

	log.Debug("saving updated configuration")
	if err := c.config.SaveConfig(ctx, cfg); err != nil {
		log.Error("failed to save configuration", logger.Error(err))
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	log.Info("configuration updated successfully", logger.String("key", opts.Key))
	fmt.Fprintf(opts.Stdout, "Set %s to %s successfully.\n", opts.Key, opts.Value)

	return nil
}

// Finalize performs any necessary cleanup after the config set operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
