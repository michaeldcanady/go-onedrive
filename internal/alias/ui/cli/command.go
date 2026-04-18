package cli

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/alias/ui/cli/list"
	"github.com/michaeldcanady/go-onedrive/internal/alias/ui/cli/remove"
	"github.com/michaeldcanady/go-onedrive/internal/alias/ui/cli/set"
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
