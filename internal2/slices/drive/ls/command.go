package ls

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// CreateLsCmd constructs and returns the cobra.Command for the drive ls operation.
func CreateLsCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List directory contents",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Path = args[0]
			}
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("drive-ls")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "list subdirectories recursively")
	cmd.Flags().StringVarP(&opts.Format, "format", "f", "short", "output format (short, long, json, yaml, tree)")
	cmd.Flags().StringVar(&opts.SortField, "sort", "Name", "field to sort by (Name, Size, ModifiedAt)")
	cmd.Flags().BoolVar(&opts.SortDescending, "desc", false, "sort in descending order")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "include hidden items")

	// Helper flags for specific formats
	var long bool
	cmd.Flags().BoolVarP(&long, "long", "l", false, "use long listing format")
	var tree bool
	cmd.Flags().BoolVar(&tree, "tree", false, "use tree listing format")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if long {
			opts.Format = "long"
		}
		if tree {
			opts.Format = "tree"
		}
	}

	return cmd
}
