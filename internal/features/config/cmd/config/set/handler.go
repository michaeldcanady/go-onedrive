package set

import (
	"fmt"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Options.Key == "" {
		return fmt.Errorf("key is required")
	}
	if ctx.Options.Value == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	return c.config.Set(ctx.Options.Key, ctx.Options.Value)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Printf("Successfully set %s to %s\n", ctx.Options.Key, ctx.Options.Value)
	return nil
}
