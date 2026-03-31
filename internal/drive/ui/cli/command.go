// Package drive provides the drive-related CLI commands.
package drive

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/alias"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/get"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/list"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/use"
	"github.com/spf13/cobra"
)

// CreateDriveCmd constructs and returns the cobra.Command for the 'drive' parent command.
func CreateDriveCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drive",
		Short: "Manage OneDrive drives and aliases",
	}

	cmd.AddCommand(
		list.CreateListCmd(container),
		use.CreateUseCmd(container),
		get.CreateGetCmd(container),
		alias.CreateAliasCmd(container),
	)

	return cmd
}
