package get

import (
	"github.com/michaeldcanady/go-onedrive/internal/config/ui/cli/shared"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateGetCmd constructs and returns the cobra.Command for the config get operation.
func CreateGetCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "get [key]",
		Short:             "Display configuration settings",
		Long:              `Display all configuration settings or a specific setting by key for the active profile.`,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: shared.ConfigKeyCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]
			opts.Stdout = cmd.OutOrStdout()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(
				container.Config(),
				container.Logger(),
			).Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
