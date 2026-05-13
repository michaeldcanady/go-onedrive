package get

import (
	"fmt"
	"github.com/michaeldcanady/go-onedrive/pkg/format"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Options.Key == "" {
		return fmt.Errorf("key is required")
	}
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	val, err := c.config.Get(ctx.Options.Key)
	if err != nil {
		return err
	}

	f := c.formatter.Get(format.Format(ctx.Options.Format))
	return f.Format(ctx.Options.Stdout, val)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
