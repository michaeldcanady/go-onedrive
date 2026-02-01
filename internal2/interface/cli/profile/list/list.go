package list

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := container.Profile().List()
			if err != nil {
				return err
			}

			for _, p := range profiles {
				cmd.Printf("%s\n", p.Name)
			}

			return nil
		},
	}
}
