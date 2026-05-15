package mkdir

import (
	"fmt"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	return c.fS.Mkdir(ctx.Ctx, ctx.Options.Path)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Printf("Created directory: %s\n", ctx.Options.Path)
	return nil
}
