package touch

import (
	"bytes"
	"fmt"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Options.Path == "" {
		return fmt.Errorf("path is required")
	}
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	// For touch, we'll just write an empty buffer if it doesn't exist
	// Or ideally we would have an update timestamp RPC.
	// For now, let's just write empty content.
	return c.fS.Write(ctx.Ctx, ctx.Options.Path, bytes.NewReader(nil))
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Printf("Touched: %s\n", ctx.Options.Path)
	return nil
}
