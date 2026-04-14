package ls

import (
	"fmt"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateLsCmd constructs and returns the cobra.Command for the ls operation.
func CreateLsCmd(container di.Container) *cobra.Command {
	var opts Options
	var format string

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

			// Resolve URI using the factory
			uri, err := container.URIFactory().FromString(opts.Path)
			if err != nil {
				return err
			}
			opts.URI = uri

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l, _ := container.Logger().CreateLogger("ls")
			handler := NewHandler(container.FS(), l)
			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "o", "short", "Output format (short, long, json, yaml, tree)")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "List items recursively")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Show hidden items")
	cmd.Flags().StringSliceVar(&opts.SortFields, "sort", []string{"name"}, "Sort items by field (name, size, modified)")
	cmd.Flags().BoolVar(&opts.SortDescending, "desc", false, "Sort in descending order")

	return cmd
}
