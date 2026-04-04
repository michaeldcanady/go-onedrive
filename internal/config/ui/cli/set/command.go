package set

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateSetCmd constructs and returns the cobra.Command for the config set operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update configuration settings",
		Long:  `Update a specific configuration setting for the active profile.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]
			opts.Value = args[1]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			handler := NewHandler(
				container.Config(),
				container.Profile(),
				container.Logger(),
			)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
