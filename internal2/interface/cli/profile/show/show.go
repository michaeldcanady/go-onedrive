package show

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

func CreateShowCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show details for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			p, err := container.Profile().Get(cmd.Context(), name)
			if err != nil {
				return err
			}

			cmd.Printf("Name: %s\n", p.Name)
			cmd.Printf("Path: %s\n", p.Path)
			if p.ConfigurationPath != "" {
				cmd.Printf("Config: %s\n", p.ConfigurationPath)
			}

			return nil
		},
	}
}
