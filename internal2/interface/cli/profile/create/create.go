package create

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

func CreateCreateCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			p, err := container.Profile().Create(name)
			if err != nil {
				return err
			}

			cmd.Printf("Created profile %q at %s\n", p.Name, p.Path)
			return nil
		},
	}
}
