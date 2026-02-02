package drive

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

// CreateDriveCmd constructs the `drive` parent command, which groups all
// drive-related subcommands.
func CreateDriveCmd(container di.Container) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "drive <subcommand>",
		Short: "Manage OneDrive drives",
	}

	// Subcommands
	//authCmd.AddCommand()

	return authCmd
}
