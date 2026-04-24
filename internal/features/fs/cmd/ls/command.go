package ls

import (
	"context"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"

	"github.com/spf13/cobra"
)

// CreateLsCmd constructs and returns the cobra.Command for the ls operation.
func CreateLsCmd(container di.Container) *cobra.Command {
	ctx := &CommandContext{
		Options: Options{},
	}
	var format string

	l, _ := container.Logger().CreateLogger("ls")
	handler := NewCommand(container.FS(), container.URIFactory(), formatting.NewFormatterFactory(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "ls <path>",
		Short:   "List items in a directory",
		Long:    "List the items in a specified directory in OneDrive or the local filesystem.",
		Args:    cobra.MaximumNArgs(1),
		Handler: handler,
		Options: ctx,
		CtxFunc: func(c context.Context, cc *CommandContext) *CommandContext {
			cc.Ctx = c
			return cc
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				ctx.Options.Path = args[0]; ctx.Options.Stdout = cmd.OutOrStdout()
			}
			ctx.Options.Format = formatting.NewFormat(format)
			return nil
		},
	})

	cmd.ValidArgsFunction = cli.ProviderPathCompletion(container)
	cmd.Flags().StringVarP(&format, "format", "o", "short", "Output format (short, long, json, yaml, tree)")
	cmd.Flags().BoolVarP(&ctx.Options.Recursive, "recursive", "r", false, "List items recursively")
	cmd.Flags().BoolVarP(&ctx.Options.All, "all", "a", false, "Show hidden items")
	cmd.Flags().StringSliceVar(&ctx.Options.SortFields, "sort", []string{"name"}, "Sort items by field (name, size, modified)")
	cmd.Flags().BoolVar(&ctx.Options.SortDescending, "desc", false, "Sort in descending order")

	return cmd
}
