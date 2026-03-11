package drive

import (
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive/get"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive/list"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive/use"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

// CreateDriveCmd constructs the `drive` parent command, which groups all
// drive-related subcommands.
func CreateDriveCmd(container didomain.Container) *cobra.Command {
	driveCmd := &cobra.Command{
		Use:   "drive <subcommand>",
		Short: "Manage OneDrive drives",
	}

	// Subcommands
	driveCmd.AddCommand(
		list.CreateListCmd(container),
		use.CreateUseCmd(container),
		alias.CreateAliasCmd(container),
		get.CreateGetCmd(container),
	)

	return driveCmd
}
