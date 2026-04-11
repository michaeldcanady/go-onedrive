package cat

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

// CreateCatCmd constructs and returns the cobra.Command for the drive cat operation.
func CreateCatCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "cat <path>",
		Short:             "Display file contents",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			uri, err := fs.ParseURI(args[0])
			if err != nil {
				return coreerrors.NewInvalidInput(
					err,
					fmt.Sprintf("invalid path '%s'", args[0]),
					"Check the path format and ensure no illegal characters are used.",
				)
			}
			opts.Path = uri

			if err := fs.ValidatePathSyntax(uri.Path); err != nil {
				switch err.(type) {
				case *coreerrors.TrailingSlashError:
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid path '%s' due to trailing slash", opts.Path.String()),
						"Remove the trailing slash from the path",
					)
				default:
					return err
				}
			}

			if names, err := container.ProviderRegistry().RegisteredNames(); err != nil {
				return coreerrors.NewAppError(
					coreerrors.CodeUnknown,
					errors.New("failed to check registered providers"),
					"An unexpected error occurred while retrieving registered providers",
					"Try again, and if the problem persists, check the application logs for more details",
				)
			} else if !slices.Contains(names, uri.Provider) {
				return coreerrors.NewInvalidInput(
					errors.New("unknown provider"),
					fmt.Sprintf("unknown provider prefix '%s'", uri.Provider),
					"Ensure the provider prefix is correct and corresponds to a registered provider",
				)
			}

			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.OutOrStderr()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.FS(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
