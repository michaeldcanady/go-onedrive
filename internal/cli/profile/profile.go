package profile

import (
	"github.com/michaeldcanady/go-onedrive/internal/cli/profile/create"
	"github.com/michaeldcanady/go-onedrive/internal/cli/profile/current"
	"github.com/michaeldcanady/go-onedrive/internal/cli/profile/delete"
	"github.com/michaeldcanady/go-onedrive/internal/cli/profile/list"
	"github.com/michaeldcanady/go-onedrive/internal/cli/profile/show"
	"github.com/michaeldcanady/go-onedrive/internal/cli/profile/use"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

func CreateProfileCmd(container didomain.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile [subcommand]",
		Short: "Manage OneDrive CLI profiles",
		Long:  "Create, list, delete, and inspect OneDrive CLI profiles.",
	}

	cmd.AddCommand(
		list.CreateListCmd(container),
		create.CreateCreateCmd(container),
		delete.CreateDeleteCmd(container),
		show.CreateShowCmd(container),
		use.CreateUseCmd(container),
		current.CreateCurrentCmd(container),
	)

	return cmd
}
