package get

import (
	"context"

	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

func CreateGetCmd(container di.Container) *cobra.Command {
	ctx := &CommandContext{
		Options: &Options{},
	}

	l, _ := container.Logger().CreateLogger("config-get")
	handler := NewCommand(container.Config(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "get [key]",
		Short:   "Get configuration",
		Args:    cobra.ExactArgs(1),
		Handler: handler,
		Options: ctx,
		CtxFunc: func(c context.Context, cc *CommandContext) *CommandContext {
			cc.Ctx = c
			return cc
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx.Options.Key = args[0]
			ctx.Options.Stdout = cmd.OutOrStdout()
			ctx.Options.Stderr = cmd.ErrOrStderr()
			return nil
		},
	})

	return cmd
}
