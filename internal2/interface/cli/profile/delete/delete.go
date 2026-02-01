package delete

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

func CreateDeleteCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if err := container.Profile().Delete(name); err != nil {
				return err
			}

			cmd.Printf("Deleted profile %q\n", name)
			return nil
		},
	}
}
