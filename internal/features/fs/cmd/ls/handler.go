package ls

import (
	"github.com/michaeldcanady/go-onedrive/pkg/format"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Options.Path == "" {
		ctx.Options.Path = "."
	}
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	nodes, err := c.fS.List(ctx.Ctx, ctx.Options.Path)
	if err != nil {
		return err
	}

	names := make([]string, len(nodes))
	for i, n := range nodes {
		names[i] = n.Name
	}

	f := c.formatter.Get(format.Format(ctx.Options.Format))
	return f.Format(ctx.Options.Stdout, names)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
