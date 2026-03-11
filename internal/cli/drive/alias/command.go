package alias

import (
	"github.com/michaeldcanady/go-onedrive/internal/cli/drive/alias/list"
	"github.com/michaeldcanady/go-onedrive/internal/cli/drive/alias/remove"
	"github.com/michaeldcanady/go-onedrive/internal/cli/drive/alias/set"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

func CreateAliasCmd(container didomain.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage drive aliases",
	}

	cmd.AddCommand(
		set.CreateSetCmd(container),
		remove.CreateRemoveCmd(container),
		list.CreateListCmd(container),
	)

	return cmd
}
