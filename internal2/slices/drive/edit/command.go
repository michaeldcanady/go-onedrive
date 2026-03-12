package edit

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
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
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

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
