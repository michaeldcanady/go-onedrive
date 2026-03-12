package create

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// CreateCreateCmd constructs and returns the cobra.Command for the profile creation operation.
func CreateCreateCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new configuration profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("profile-create")
			handler := NewHandler(container.Profile(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
