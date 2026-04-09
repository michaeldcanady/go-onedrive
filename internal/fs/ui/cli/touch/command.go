package touch

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
			if err := fs.ValidatePathSyntax(opts.Path); err != nil {
				switch err.(type) {
				case *fs.TrailingSlashError:
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid path '%s' due to trailing slash", opts.Path),
						"Remove the trailing slash from the path",
					)
				case *fs.IllegalCharacterError:
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid path '%s' due to illegal characters", opts.Path),
						"Remove the illegal characters from the path",
					)
				default:
					return err
				}
			}

			// 2. Provider check (only if a provider prefix is explicitly given)
			provider, _, found := fs.SplitProviderPath(opts.Path)
			if found {
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

	return cmd
}
