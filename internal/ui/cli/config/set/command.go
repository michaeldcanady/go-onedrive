package set

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateSetCmd constructs and returns the cobra.Command for the config set operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	var opts Options

	l, _ := container.Logger().CreateLogger("config-set")
	handler := NewCommand(
		container.Config(),
		l,
	)

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update a configuration setting",
		Long:  `Update a specific configuration setting for the active profile.`,
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]
			opts.Value = args[1]
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
