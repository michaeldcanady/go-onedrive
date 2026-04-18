package add

import (
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

func CreateAddCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "add <path> <type> <identity_id>",
		Short: "Add a mount point",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.Mounts().AddMount(cmd.Context(), config.MountConfig{
				Path:       args[0],
				Type:       args[1],
				IdentityID: args[2],
			})
		},
	}
}
