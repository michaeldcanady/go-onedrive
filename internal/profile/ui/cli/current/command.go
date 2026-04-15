package current

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateCurrentCmd constructs and returns the cobra.Command for the profile current operation.
func CreateCurrentCmd(container di.Container) *cobra.Command {
	var opts Options

	l, _ := container.Logger().CreateLogger("profile-current")
	handler := NewCommand(container.Profile(), l)

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the active profile",
		Long:  "Display the name of the currently active profile.",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()

			return handler.Validate(cmd.Context(), &opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return handler.Execute(cmd.Context(), opts)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return handler.Finalize(cmd.Context(), opts)
		},
	}

	return cmd
}
