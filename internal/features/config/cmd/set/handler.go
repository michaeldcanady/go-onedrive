package set

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
)

// Command executes the config set operation.
type Command struct {
	config config.Service
	log    logger.Logger
}

var illegalChars = []string{
	" ",
	"\n",
	"\r",
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
	cleanKey := strings.TrimSpace(ctx.Options.Key)
	if cleanKey == "" {
		return errors.New("key is empty")
	}

	for _, illegalChar := range illegalChars {
		if strings.Contains(cleanKey, illegalChar) {
			return fmt.Errorf("key contains illegal char %s", illegalChar)
		}
	}

	cleanValue := strings.TrimSpace(ctx.Options.Value)
	if cleanValue == "" {
		return errors.New("value is empty")
	}

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
