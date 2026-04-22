// Package profile provides the profile-related CLI commands.
package profile

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/profile/create"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/profile/current"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/profile/delete"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/profile/list"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/profile/use"
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
		create.CreateCreateCmd(container),
		delete.CreateDeleteCmd(container),
		use.CreateUseCmd(container),
	)

	return cmd
}
