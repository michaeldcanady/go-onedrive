package stat

import (
	"fmt"
	"time"
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
	n, err := c.fS.Stat(ctx.Ctx, ctx.Options.Path)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", n.Name)
	fmt.Printf("Path: %s\n", n.Path)
	fmt.Printf("Size: %d bytes\n", n.Size)
	fmt.Printf("Modified: %v\n", time.Unix(n.ModifiedAt, 0))
	return nil
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
