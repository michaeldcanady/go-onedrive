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
	opts := Options{}
	var format string

	l, _ := container.Logger().CreateLogger("ls")
	handler := NewCommand(container.FS(), container.URIFactory(), formatting.NewFormatterFactory(), l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "ls <path>",
		Short:   "List items in a directory",
		Long:    "List the items in a specified directory in OneDrive or the local filesystem.",
		Args:    cobra.MaximumNArgs(1),
		Handler: handler,
		Options: &CommandContext{Options: opts},
		CtxFunc: func(ctx context.Context, c *CommandContext) *CommandContext {
			c.Ctx = ctx
			return c
		},
	})

	cmd.ValidArgsFunction = cli.ProviderPathCompletion(container)
	cmd.Flags().StringVarP(&format, "format", "o", "short", "Output format (short, long, json, yaml, tree)")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "List items recursively")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Show hidden items")
	cmd.Flags().StringSliceVar(&opts.SortFields, "sort", []string{"name"}, "Sort items by field (name, size, modified)")
	cmd.Flags().BoolVar(&opts.SortDescending, "desc", false, "Sort in descending order")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			opts.Path = args[0]
		}
	}

	return cmd
}
