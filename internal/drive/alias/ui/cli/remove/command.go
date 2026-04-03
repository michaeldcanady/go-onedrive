package remove

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/alias/shared"
	"github.com/spf13/cobra"
)

// CreateRemoveCmd constructs and returns the cobra.Command for the 'drive alias remove' operation.
func CreateRemoveCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:               "remove <alias>",
		Short:             "Remove a drive alias",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate that the alias exists before attempting to remove it
			alias := args[0]
			driveID, err := container.State().GetDriveAlias(alias)
			if err != nil {
				return err
			}
			if driveID == "" {
				return fmt.Errorf("alias '%s' not found", alias)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Alias:  args[0],
				Stdout: cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("alias-remove")
			return NewHandler(container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
