package cat

import (
	"io"
	"os"
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
	reader, err := c.fS.Read(ctx.Ctx, ctx.Options.Path)
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	return err
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
