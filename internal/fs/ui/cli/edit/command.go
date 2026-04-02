package edit

import (
	"fmt"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/di"
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
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

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
			log, err := container.Logger().CreateLogger("edit")
			if err != nil {
				return err
			}

			handler := NewHandler(container.FS(), container.Editor(), log)
			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite even if conflicts are detected")

	return cmd
}
