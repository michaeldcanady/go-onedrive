package mount

import (
	"github.com/michaeldcanady/go-onedrive/internal/features/di"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/mount/add"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/mount/list"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/mount/remove"

	"github.com/spf13/cobra"
)

// CreateMountCmd constructs the 'mount' parent command.
func CreateMountCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mount",
		Short: "Manage virtual filesystem mount points",
	}

	cmd.AddCommand(
		list.CreateListCmd(container),
		add.CreateAddCmd(container),
		remove.CreateRemoveCmd(container),
	)

	return cmd
}
