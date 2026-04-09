package ls

import (
	"errors"
	"fmt"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
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
			opts.Path = args[0]
			// should allow / for directories
			if contains, err := fs.ContainsIllegalChars(opts.Path); contains {
				return coreerrors.NewInvalidInput(
					err,
					fmt.Sprintf("invalid path '%s' due to illegal characters", opts.Path),
					"Remove the illegal characters from the path",
				)
			}

			if provider, _, found := fs.SplitProviderPath(opts.Path); found {
				if names, err := container.ProviderRegistry().RegisteredNames(); err != nil {
					return coreerrors.NewAppError(
						coreerrors.CodeUnknown,
						errors.New("failed to check registered providers"),
						"An unexpected error occurred while retrieving registered providers",
						"Try again, and if the problem persists, check the application logs for more details",
					)
				} else if !slices.Contains(names, provider) {
					return coreerrors.NewInvalidInput(
						errors.New("unknown provider"),
						"Unknown provider prefix",
						"Ensure the provider prefix is correct and corresponds to a registered provider",
					)
				}
			}

			opts.Stdout = cmd.OutOrStdout()
			opts.Format = formatting.NewFormat(format)

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.FS(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "o", "short", "Output format (short, long, json, yaml, tree)")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "List items recursively")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Show hidden items")
	cmd.Flags().StringVar(&opts.SortField, "sort", "name", "Sort items by field (name, size, modified)")
	cmd.Flags().BoolVar(&opts.SortDescending, "desc", false, "Sort in descending order")

	return cmd
}
