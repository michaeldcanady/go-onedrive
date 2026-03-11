package edit

import (
	"fmt"

	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "edit"
)

// CreateEditCmd constructs and returns the cobra.Command for the edit operation.
// It initializes flags and sets up the execution logic using EditCmd.
func CreateEditCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <PATH>", commandName),
		Short: "Edit a OneDrive file in your local editor",
		Long: `You can edit a OneDrive file using your favorite local text editor. This
command downloads the file to a temporary location, opens it in your editor,
and automatically uploads the changes back to OneDrive when you save and
close the file.`,
		Example: `  # Edit a text file in your default editor
  odc edit /Notes/ideas.txt

  # Force an upload even if the file has changed on OneDrive
  odc edit /Projects/config.yaml --force`,
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}

			opts.Stdin = cmd.InOrStdin()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			editCmd := NewEditCmd(c)
			return editCmd.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite the file even if it changed in the cloud")

	return cmd
}
