package mount

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount/cmd/add"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount/cmd/list"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount/cmd/remove"

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
