package set

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

func CreateSetCmd(container di.Container) *cobra.Command {
	opts := NewOptions()

	l, _ := container.Logger().CreateLogger("mount-add")
	handler := NewCommand(container.Config(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "add <path> <type> <identity_id>",
		Short:   "Add a mount point",
		Args:    cobra.ExactArgs(3),
		Handler: handler,
		Options: NewCommandContext(context.Background(), opts),
		CtxFunc: func(ctx context.Context, o *CommandContext) *CommandContext {
			o.Ctx = ctx
			return o
		},
	})

	return cmd
}
