package mkdir

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

// CreateMkdirCmd constructs and returns the cobra.Command for the drive mkdir operation.
func CreateMkdirCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "mkdir <path>",
		Short:             "Create a new directory",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			if err := fs.ValidatePathSyntax(opts.Path); err != nil {
				switch err.(type) {
				case *fs.TrailingSlashError:
					break
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

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.FS(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
