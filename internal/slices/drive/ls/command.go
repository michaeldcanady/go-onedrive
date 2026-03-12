package ls

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateLsCmd constructs and returns the cobra.Command for the ls operation.
func CreateLsCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List items in a directory",
		Long:  "List the items in a specified directory in OneDrive or the local filesystem.",
		Args:  cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// This is a basic implementation. A real one would need to:
			// 1. Initialize a lightweight container.
			// 2. Resolve the path up to `toComplete`.
			// 3. List children of that path.
			// 4. Return their names.
			return nil, cobra.ShellCompDirectiveDefault
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Path = args[0]
			}

			l, _ := container.Logger().CreateLogger("ls")
			handler := NewHandler(container.FS(), l)
			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Format, "format", "o", "short", "Output format (short, long, json, yaml, tree)")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "List items recursively")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Show hidden items")
	cmd.Flags().StringVar(&opts.SortField, "sort", "name", "Sort items by field (name, size, modified)")
	cmd.Flags().BoolVar(&opts.SortDescending, "desc", false, "Sort in descending order")

	return cmd
}
