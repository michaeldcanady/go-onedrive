package add

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
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
	m := &mount.Mount{
		Path:       ctx.Options.Path,
		Type:       ctx.Options.Type,
		IdentityID: ctx.Options.IdentityId,
		Options:    c.parseOptions(ctx.Options.Option),
	}

	return c.mounts.Add(ctx.Ctx, m)
}

func (c *Command) parseOptions(opts []string) map[string]string {
	options := make(map[string]string)
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) == 2 {
			options[parts[0]] = parts[1]
		}
	}
	return options
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Printf("Successfully added mount point: %s\n", ctx.Options.Path)
	return nil
}
