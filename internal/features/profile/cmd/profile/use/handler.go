package use

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
	return c.profile.SetCurrent(ctx.Options.Name)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Printf("Now using profile: %s\n", ctx.Options.Name)
	return nil
}
