package touch

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"

	"github.com/spf13/cobra"
)

// CreateTouchCmd constructs and returns the cobra.Command for the drive touch operation.
func CreateTouchCmd(container di.Container) *cobra.Command {
	ctx := &CommandContext{
		Options: Options{},
	}

	l, _ := container.Logger().CreateLogger("touch")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "touch <path>",
		Short:   "Create a new empty file",
		Args:    cobra.ExactArgs(1),
		Handler: handler,
		Options: ctx,
		CtxFunc: func(c context.Context, cc *CommandContext) *CommandContext {
			cc.Ctx = c
			return cc
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx.Options.Path = args[0]
			ctx.Options.Stdout = cmd.OutOrStdout()
			return nil
		},
	})

	cmd.ValidArgsFunction = cli.ProviderPathCompletion(container)

	return cmd
}
