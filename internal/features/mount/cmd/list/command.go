package list

import (
	"context"
	"github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	"github.com/spf13/cobra"
)

func CreateListCmd(container di.Container) *cobra.Command {
	opts := NewOptions()

	l, _ := container.Logger().CreateLogger("mount-list")
	handler := NewCommand(container.Mounts(), formatting.NewFormatterFactory(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "list",
		Short:   "List mount points",
		Handler: handler,
		Options: NewCommandContext(context.Background(), opts),
		CtxFunc: func(ctx context.Context, o *CommandContext) *CommandContext {
			o.Ctx = ctx
			return o
		},
	})

	formatStr := string(opts.Format)
	cmd.Flags().StringVarP(&formatStr, "format", "f", "short", "output format (short, long, json, yaml, tree, table)")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		opts.Format = formatStr
	}

	return cmd
}
