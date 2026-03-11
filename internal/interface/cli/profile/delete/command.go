package delete

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "delete"
	loggerID    = "cli"
)

// CreateDeleteCmd constructs and returns the cobra.Command for the delete operation.
func CreateDeleteCmd(container didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "delete <NAME>",
		Short: "Delete a profile",
		Long: `You can delete a OneDrive profile that you no longer need. This removes the
associated configuration and authentication tokens from your local machine.`,
		Example: `  # Delete the profile named 'old-account'
  odc profile delete old-account

  # Forcefully delete a profile without asking for confirmation
  odc profile delete temporary-profile --force`,
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			return NewDeleteCmd(container).Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force deletion without confirmation")

	return cmd
}
