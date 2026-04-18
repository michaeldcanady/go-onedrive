package ls

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"

	"github.com/spf13/cobra"
)

// CreateLsCmd constructs and returns the cobra.Command for the ls operation.
func CreateLsCmd(container di.Container) *cobra.Command {
	var opts Options
	var format string
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("ls")
	handler := NewCommand(container.FS(), container.URIFactory(), formatting.NewFormatterFactory(), l)

	cmd := &cobra.Command{
		Use:               "ls <path>",
		Short:             "List items in a directory",
		Long:              "List the items in a specified directory in OneDrive or the local filesystem.",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Path = args[0]
			}
			opts.Stdout = cmd.OutOrStdout()
			opts.Format = formatting.NewFormat(format)

			c = &CommandContext{
				Ctx:     cmd.Context(),
				Options: opts,
			}

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "o", "short", "Output format (short, long, json, yaml, tree)")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "List items recursively")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Show hidden items")
	cmd.Flags().StringSliceVar(&opts.SortFields, "sort", []string{"name"}, "Sort items by field (name, size, modified)")
	cmd.Flags().BoolVar(&opts.SortDescending, "desc", false, "Sort in descending order")

	return cmd
}
