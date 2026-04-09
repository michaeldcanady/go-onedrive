package edit

import (
	"errors"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateEditCmd constructs and returns the cobra.Command for the edit operation.
func CreateEditCmd(container di.Container) *cobra.Command {
	opts := NewOptions()

	cmd := &cobra.Command{
		Use:   "edit <path>",
		Short: "Edit a file in an external editor",
		Long: `Download a file to a temporary location, open it in your preferred
editor ($VISUAL, $EDITOR, or system defaults), and upload the changes back
to OneDrive.`,
		Example: `  # Edit a file in your OneDrive root
  odc drive edit document.txt

  # Force overwrite even if changes exist on server
  odc drive edit -f document.txt`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			if err := fs.ValidatePathSyntax(opts.Path); err != nil {
				return err
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
			opts.Stderr = cmd.ErrOrStderr()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.FS(), container.Editor(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite even if conflicts are detected")

	return cmd
}
