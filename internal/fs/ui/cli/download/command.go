package download

import (
	"fmt"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateDownloadCmd constructs and returns the cobra.Command for the drive download operation.
func CreateDownloadCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "download <remote_path> <local_path>",
		Short:             "Download files and directories from OneDrive",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			opts.Stdout = cmd.OutOrStdout()

			// 1. Syntactic check for source path
			if err := fs.ValidatePathSyntax(opts.Source); err != nil {
				return fmt.Errorf("invalid source path syntax: %w", err)
			}
			// 2. Provider check for source path (only if a provider prefix is explicitly given)
			provider, _, found := fs.SplitProviderPath(opts.Source)
			if found {
				names, err := container.ProviderRegistry().RegisteredNames()
				if err != nil {
					return fmt.Errorf("failed to check registered providers: %w", err)
				}
				if !slices.Contains(names, provider) {
					return fmt.Errorf("unknown provider '%s'; valid providers are: %v", provider, names)
				}
			}

			// 3. Syntactic check for destination path
			if err := fs.ValidatePathSyntax(opts.Destination); err != nil {
				return fmt.Errorf("invalid destination path syntax: %w", err)
			}
			// 4. Provider check for destination path (only if a provider prefix is explicitly given)
			destProvider, _, destFound := fs.SplitProviderPath(opts.Destination)
			if destFound {
				names, err := container.ProviderRegistry().RegisteredNames()
				if err != nil {
					return fmt.Errorf("failed to check registered providers: %w", err)
				}
				if !slices.Contains(names, destProvider) {
					return fmt.Errorf("unknown provider '%s'; valid providers are: %v", destProvider, names)
				}
			}

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l, _ := container.Logger().CreateLogger("drive-download")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "download directories recursively")

	return cmd
}
