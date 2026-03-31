package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for listing profiles.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()

			l, _ := container.Logger().CreateLogger("profile-list")
			handler := NewHandler(container.Profile(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
