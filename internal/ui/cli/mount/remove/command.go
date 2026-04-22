package remove

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

func CreateRemoveCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <path>",
		Short: "Remove a mount point",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.Mounts().RemoveMount(cmd.Context(), args[0])
		},
	}
}
