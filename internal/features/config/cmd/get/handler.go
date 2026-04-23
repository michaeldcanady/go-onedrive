package get

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
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

var illegalChars = []string{
	" ",
	"\n",
	"\r",
}

// Validate prepares and validates the options for the config get operation.
func (c *Command) Validate(ctx *CommandContext) error {
	cleanKey := strings.TrimSpace(ctx.Options.Key)
	if cleanKey == "" {
		return errors.New("key is empty")
	}

	for _, illegalChar := range illegalChars {
		if strings.Contains(cleanKey, illegalChar) {
			return fmt.Errorf("key contains illegal char %s", illegalChar)
		}
	}

	return ctx.Options.Validate()
}

// Execute retrieves and displays configuration settings.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Debug("fetching configuration")
	cfg, err := c.config.GetConfig(ctx.Ctx)
	if err != nil {
		log.Error("failed to retrieve configuration", logger.Error(err))
		return fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	if ctx.Options.Key != "" {
		log.Debug("retrieving specific configuration key", logger.String("key", ctx.Options.Key))
		// TODO: define a generic way to discover key values in the appropriate format?
		// For now, we'll just handle a few common keys manually until we have a generic way.
		var val interface{}
		switch ctx.Options.Key {
		case "auth.provider":
			val = cfg.Auth.Provider
		case "auth.method":
			val = cfg.Auth.Method
		case "logging.level":
			val = cfg.Logging.Level
		case "logging.format":
			val = cfg.Logging.Format
		default:
			return fmt.Errorf("configuration key not supported via CLI: %s", ctx.Options.Key)
		}
		fmt.Fprintf(ctx.Options.Stdout, "%s: %v\n", ctx.Options.Key, val)
		return nil
	}

	// TODO: this needs to be handled by the formatter
	log.Debug("displaying full configuration")
	var data []byte
	switch cfg.Logging.Format {
	case "json":
		data, _ = json.MarshalIndent(cfg, "", "  ")
	default:
		data, _ = yaml.Marshal(cfg)
	}

	fmt.Fprintln(ctx.Options.Stdout, string(data))
	return nil
}

// Finalize performs any necessary cleanup after the config get operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
