package logout

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
	if ctx.Options.Id == "" {
		return fmt.Errorf("identity ID is required")
	}
	return c.identity.Logout(ctx.Ctx, ctx.Options.Id)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	if ctx.Options.Id != "" {
		fmt.Printf("Logged out from account: %s\n", ctx.Options.Id)
	} else {
		fmt.Println("Logged out from all accounts for the active profile")
	}
	return nil
}
