package upload

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"

	"github.com/spf13/cobra"
)

// CreateUploadCmd constructs and returns the cobra.Command for the drive upload operation.
func CreateUploadCmd(container di.Container) *cobra.Command {
	ctx := &CommandContext{
		Options: Options{},
	}

	l, _ := container.Logger().CreateLogger("upload")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "upload <source> <destination>",
		Short:   "Upload files and directories",
		Args:    cobra.ExactArgs(2),
		Handler: handler,
		Options: ctx,
		CtxFunc: func(c context.Context, cc *CommandContext) *CommandContext {
			cc.Ctx = c
			return cc
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx.Options.Source = args[0]
			ctx.Options.Destination = args[1]
			ctx.Options.Stdout = cmd.OutOrStdout()
			return nil
		},
	})

	cmd.ValidArgsFunction = cli.ProviderPathCompletion(container)
	cmd.Flags().BoolVarP(&ctx.Options.Recursive, "recursive", "r", false, "upload directories recursively")

	return cmd
}
