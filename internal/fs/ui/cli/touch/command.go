package touch

import (
	"fmt"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateTouchCmd constructs and returns the cobra.Command for the drive touch operation.
func CreateTouchCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "touch <path>",
		Short:             "Create an empty file or update its timestamp",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			opts.Stdout = cmd.OutOrStdout()

			// 1. Syntactic check
			if err := fs.ValidatePathSyntax(opts.Path); err != nil {
				return fmt.Errorf("invalid path syntax: %w", err)
			}

			// 2. Provider check (only if a provider prefix is explicitly given)
			provider, _, found := fs.SplitProviderPath(opts.Path)
			if found {
				names, err := container.ProviderRegistry().RegisteredNames()
				if err != nil {
					return fmt.Errorf("failed to check registered providers: %w", err)
				}
				if !slices.Contains(names, provider) {
					return fmt.Errorf("unknown provider '%s'; valid providers are: %v", provider, names)
				}
			}

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l, _ := container.Logger().CreateLogger("drive-touch")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
