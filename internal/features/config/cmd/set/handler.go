package set

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
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
func (c *Command) Validate(ctx *CommandContext) error {
	return ctx.Options.Validate()
}

// Execute updates a configuration setting.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Debug("updating configuration setting", logger.String("key", ctx.Options.Key), logger.String("value", ctx.Options.Value))
	if err := c.config.UpdateConfig(ctx.Ctx, ctx.Options.Key, ctx.Options.Value); err != nil {
		log.Error("failed to update configuration", logger.Error(err))
		return fmt.Errorf("failed to update configuration: %w", err)
	}

	log.Info("configuration updated successfully", logger.String("key", ctx.Options.Key))
	fmt.Fprintf(ctx.Options.Stdout, "Set %s to %s successfully.\n", ctx.Options.Key, ctx.Options.Value)

	return nil
}

// Finalize performs any necessary cleanup after the config set operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
