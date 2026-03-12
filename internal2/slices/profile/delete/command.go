package delete

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// CreateDeleteCmd constructs and returns the cobra.Command for the profile deletion operation.
func CreateDeleteCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a configuration profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("profile-delete")
			handler := NewHandler(container.Profile(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
