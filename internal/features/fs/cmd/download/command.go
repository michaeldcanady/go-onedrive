package download

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/spf13/cobra"
)

// CreateDownloadCmd constructs and returns the cobra.Command for the drive download operation.
func CreateDownloadCmd(container di.Container) *cobra.Command {
	ctx := &CommandContext{
		Options: Options{},
	}

	l, _ := container.Logger().CreateLogger("download")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "download <source> <destination>",
		Short:   "Download files and directories",
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
	cmd.Flags().BoolVarP(&ctx.Options.Recursive, "recursive", "r", false, "download directories recursively")

	return cmd
}

