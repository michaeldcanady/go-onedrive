package get

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"gopkg.in/yaml.v3"
)

// Command executes the config get operation.
type Command struct {
	config config.Service
	log    logger.Logger
}

// NewCommand initializes a new instance of the config get Command.
func NewCommand(c config.Service, l logger.Logger) *Command {
	return &Command{
		config: c,
		log:    l,
	}
}

// Validate prepares and validates the options for the config get operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute retrieves and displays configuration settings.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Debug("fetching configuration")
	cfg, err := c.config.GetConfig(ctx)
	if err != nil {
		log.Error("failed to retrieve configuration", logger.Error(err))
		return fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	if opts.Key != "" {
		log.Debug("retrieving specific configuration key", logger.String("key", opts.Key))
		// For now, we'll just handle a few common keys manually until we have a generic way.
		var val interface{}
		switch opts.Key {
		case "auth.provider":
			val = cfg.Auth.Provider
		case "auth.method":
			val = cfg.Auth.Method
		case "logging.level":
			val = cfg.Logging.Level
		case "logging.format":
			val = cfg.Logging.Format
		default:
			return fmt.Errorf("configuration key not supported via CLI: %s", opts.Key)
		}
		fmt.Fprintf(opts.Stdout, "%s: %v\n", opts.Key, val)
		return nil
	}

	log.Debug("displaying full configuration")
	var data []byte
	switch cfg.Logging.Format {
	case "json":
		data, _ = json.MarshalIndent(cfg, "", "  ")
	default:
		data, _ = yaml.Marshal(cfg)
	}

	fmt.Fprintln(opts.Stdout, string(data))
	return nil
}

// Finalize performs any necessary cleanup after the config get operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
