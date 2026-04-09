package set

import (
	"github.com/michaeldcanady/go-onedrive/internal/config/ui/cli/shared"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateSetCmd constructs and returns the cobra.Command for the config set operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "set <key> <value>",
		Short:             "Update configuration settings",
		Long:              `Update a specific configuration setting for the active profile.`,
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: shared.SetCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]
			opts.Value = args[1]
			opts.Stdout = cmd.OutOrStdout()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(
				container.Config(),
				container.Profile(),
				container.Logger(),
			).Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
