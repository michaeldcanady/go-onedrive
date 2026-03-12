// Package profile provides the profile-related CLI commands.
package profile

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/profile/current"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/profile/list"
	"github.com/spf13/cobra"
)

// CreateProfileCmd constructs and returns the cobra.Command for the 'profile' parent command.
func CreateProfileCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage configuration profiles",
	}

	cmd.AddCommand(
		list.CreateListCmd(container),
		current.CreateCurrentCmd(container),
	)

	return cmd
}
