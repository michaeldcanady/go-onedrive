package current

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
	p, err := c.profile.GetCurrent()
	if err != nil {
		return err
	}
	if p == nil {
		fmt.Println("No profile currently active")
		return nil
	}
	fmt.Printf("Current profile: %s\n", p.Name)
	return nil
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
