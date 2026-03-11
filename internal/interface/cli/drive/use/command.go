package use

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "use"
	loggerID    = "cli"
)

// CreateUseCmd constructs and returns the cobra.Command for the use operation.
func CreateUseCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use <DRIVE_ID_OR_ALIAS>",
		Short: "Sets the active drive",
		Long: `You can set a specific drive as the active drive for your current session.
Subsequent commands that interact with OneDrive will use this drive unless you
specify otherwise. You can use either the drive's ID or an alias.`,
		Example: `  # Set the active drive using its ID
  odc drive use b!1234567890abcdef

  # Set the active drive using an alias
  odc drive use work-drive`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveIDOrAlias: args[0],
				Stdout:         cmd.OutOrStdout(),
				Stderr:         cmd.ErrOrStderr(),
			}

			return NewUseCmd(container).Run(cmd.Context(), opts)
		},
	}
}
