package download

import (
	"errors"
	"fmt"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
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
			if err := fs.ValidatePathSyntax(opts.Source); err != nil {
				return fmt.Errorf("invalid source path syntax: %w", err)
			}
			if provider, _, found := fs.SplitProviderPath(opts.Source); found {
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

			opts.Destination = args[1]
			// 3. Syntactic check for destination path
			if err := fs.ValidatePathSyntax(opts.Destination); err != nil {
				return fmt.Errorf("invalid destination path syntax: %w", err)
			}
			// 4. Provider check for destination path (only if a provider prefix is explicitly given)
			if provider, _, found := fs.SplitProviderPath(opts.Destination); found {
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
			opts.Stderr = cmd.ErrOrStderr()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.FS(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "download directories recursively")

	return cmd
}
