package alias

import (
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive/alias/list"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive/alias/remove"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive/alias/set"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

func CreateAliasCmd(container didomain.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage drive aliases",
		Long: `You can manage short names (aliases) for your OneDrive drive IDs. Aliases
make it easier to switch between different drives or refer to them in
commands without typing long IDs.`,
		Example: `  # List all your drive aliases
  odc drive alias list

  # Set a new alias for a drive
  odc drive alias set work b!1234567890abcdef

  # Remove an existing alias
  odc drive alias remove work`,
	}

	cmd.AddCommand(
		set.CreateSetCmd(container),
		remove.CreateRemoveCmd(container),
		list.CreateListCmd(container),
	)

	return cmd
}
