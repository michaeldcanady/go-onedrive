package alias

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/alias/list"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/alias/remove"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/alias/set"
	"github.com/spf13/cobra"
)

// CreateAliasCmd constructs and returns the cobra.Command for the 'drive alias' parent command.
func CreateAliasCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage drive aliases",
	}

	cmd.AddCommand(
		list.CreateListCmd(container),
		set.CreateSetCmd(container),
		remove.CreateRemoveCmd(container),
	)

	return cmd
}
