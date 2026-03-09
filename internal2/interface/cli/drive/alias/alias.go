package alias

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/drive/alias/list"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/drive/alias/remove"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/drive/alias/set"
	"github.com/spf13/cobra"
)

func CreateAliasCmd(container di.Container) *cobra.Command {
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
