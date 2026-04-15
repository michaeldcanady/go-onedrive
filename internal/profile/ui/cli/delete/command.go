package delete

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateDeleteCmd constructs and returns the cobra.Command for the profile delete operation.
func CreateDeleteCmd(container di.Container) *cobra.Command {
	var opts Options

	l, _ := container.Logger().CreateLogger("profile-delete")
	handler := NewCommand(container.Profile(), l)

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Long:  "Permanently delete a profile and its associated configuration and authentication tokens.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
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
