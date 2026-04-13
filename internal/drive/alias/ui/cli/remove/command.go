package remove

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateRemoveCmd constructs and returns the cobra.Command for the 'drive alias remove' operation.
func CreateRemoveCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:               "remove <alias>",
		Short:             "Remove a drive alias",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate that the alias exists before attempting to remove it
			alias := args[0]
			if _, err := container.Alias().GetDriveIDByAlias(alias); err != nil {
				return err
			}

			opts.Alias = args[0]
			opts.Stdout = cmd.OutOrStdout()

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
