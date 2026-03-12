package use

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateUseCmd constructs and returns the cobra.Command for the profile switch operation.
func CreateUseCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "use [name]",
		Short: "Switch the active configuration profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("profile-use")
			handler := NewHandler(container.Profile(), container.State(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
