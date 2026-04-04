package get

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateGetCmd constructs and returns the cobra.Command for the config get operation.
func CreateGetCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Display configuration settings",
		Long:  `Display all configuration settings or a specific setting by key for the active profile.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Key = args[0]
			}
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			handler := NewHandler(
				container.Config(),
				container.Logger(),
			)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
