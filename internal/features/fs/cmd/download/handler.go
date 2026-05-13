package download

import (
	"fmt"
	"io"
	"os"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Options.Source == "" {
		return fmt.Errorf("source path is required")
	}
	if ctx.Options.Destination == "" {
		return fmt.Errorf("destination path is required")
	}
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	reader, err := c.fS.Read(ctx.Ctx, ctx.Options.Source)
	if err != nil {
		return err
	}
	defer reader.Close()

	f, err := os.Create(ctx.Options.Destination)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Printf("Downloaded %s to %s\n", ctx.Options.Source, ctx.Options.Destination)
	return nil
}
