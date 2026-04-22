// Package cli provides the drive-related CLI commands.
package cli

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/drive/get"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/drive/list"
	"github.com/spf13/cobra"
)

// CreateDriveCmd constructs and returns the cobra.Command for the 'backend-discovery' parent command.
func CreateDriveCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drive",
		Short: "Discover available drives",
	}

	cmd.AddCommand(
		list.CreateListCmd(container),
		get.CreateGetCmd(container),
	)

	return cmd
}
