package mv

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

// CreateMvCmd constructs and returns the cobra.Command for the drive mv operation.
func CreateMvCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "mv <source> <destination>",
		Short:             "Move files and directories",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			if err := fs.ValidatePathSyntax(opts.Source); err != nil {
				switch err.(type) {
				case *fs.TrailingSlashError:
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid source path '%s' due to trailing slash", opts.Source),
						"Remove the trailing slash from the source path",
					)
				case *fs.IllegalCharacterError:
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid source path '%s' due to illegal characters", opts.Source),
						"Remove the illegal characters from the source path",
					)
				default:
					return err
				}
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
			if err := fs.ValidatePathSyntax(opts.Destination); err != nil {
				switch err.(type) {
				case *fs.TrailingSlashError:
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid destination path '%s' due to trailing slash", opts.Destination),
						"Remove the trailing slash from the destination path",
					)
				case *fs.IllegalCharacterError:
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid destination path '%s' due to illegal characters", opts.Destination),
						"Remove the illegal characters from the destination path",
					)
				default:
					return err
				}
			}

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

	// TODO: add support for recursive moves (e.g. moving a directory with contents)

	return cmd
}
