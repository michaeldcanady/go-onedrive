package get

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
	drive, err := c.drive.Get(ctx.Ctx, ctx.Options.DriveRef)
	if err != nil {
		return err
	}
	if drive == nil {
		return fmt.Errorf("drive not found: %s", ctx.Options.DriveRef)
	}

	fmt.Printf("Drive ID:   %s\n", drive.ID)
	fmt.Printf("Name:       %s\n", drive.Name)
	fmt.Printf("Type:       %s\n", drive.Type)
	fmt.Printf("Identity:   %s\n", drive.IdentityID)

	return nil
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
