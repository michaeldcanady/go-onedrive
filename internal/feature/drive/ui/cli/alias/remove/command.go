package remove

import (
	"github.com/michaeldcanady/go-onedrive/internal/feature/di"
	"github.com/spf13/cobra"
)

// CreateRemoveCmd constructs and returns the cobra.Command for the 'drive alias remove' operation.
func CreateRemoveCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove a drive alias",
		Args:  cobra.ExactArgs(1),
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
