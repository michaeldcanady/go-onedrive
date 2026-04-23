package remove

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

func CreateRemoveCmd(container di.Container) *cobra.Command {
	opts := NewOptions()

	l, _ := container.Logger().CreateLogger("mount-remove")
	handler := NewCommand(container.Mounts(), container.URIFactory(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "remove <path>",
		Short:   "Remove a mount point",
		Handler: handler,
		Options: NewCommandContext(context.Background(), opts),
		CtxFunc: func(ctx context.Context, o *CommandContext) *CommandContext {
			o.Ctx = ctx
			return o
		},
	})

	return cmd
}
