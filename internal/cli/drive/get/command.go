package get

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "get"
	loggerID    = "cli"
)

// CreateGetCmd constructs and returns the cobra.Command for the get operation.
func CreateGetCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id|alias>",
		Short: "Get information of named drive",
		Long: `You can retrieve detailed information about a specific OneDrive drive by
providing its ID or an alias you've previously set. This information includes
the drive's name, type, and quota details.`,
		Example: `  # Get information for a drive using its ID
  odc drive get b!1234567890abcdef

  # Get information for a drive using an alias
  odc drive get personal-drive`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveIDOrAlias: args[0],
				Stdout:         cmd.OutOrStdout(),
				Stderr:         cmd.ErrOrStderr(),
			}

			return NewGetCmd(container).Run(cmd.Context(), opts)
		},
	}
}
