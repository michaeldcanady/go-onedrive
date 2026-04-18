package list

import (
	"fmt"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List mount points",
		RunE: func(cmd *cobra.Command, args []string) error {
			mounts, err := container.Mounts().ListMounts(cmd.Context())
			if err != nil {
				return err
			}
			for _, m := range mounts {
				fmt.Printf("Path: %s, Type: %s, Identity: %s\n", m.Path, m.Type, m.IdentityID)
			}
			return nil
		},
	}
}
