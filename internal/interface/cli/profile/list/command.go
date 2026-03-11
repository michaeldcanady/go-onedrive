package list

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "list"
	loggerID    = "cli"
)

// CreateListCmd constructs and returns the cobra.Command for the list operation.
func CreateListCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		Long: `You can list all the OneDrive profiles you've created. This helps you see
which accounts and configurations are available to use.`,
		Example: `  # Display a list of all available profiles
  odc profile list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewListCmd(container).Run(cmd.Context(), opts)
		},
	}
}
