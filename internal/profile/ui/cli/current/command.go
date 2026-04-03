package current

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateCurrentCmd constructs and returns the cobra.Command for showing the current profile.
func CreateCurrentCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Display the name of the currently active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()

			l, _ := container.Logger().CreateLogger("profile-current")
			handler := NewHandler(container.Profile(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
