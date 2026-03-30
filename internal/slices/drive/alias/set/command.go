package set

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateSetCmd constructs and returns the cobra.Command for the 'drive alias set' operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "set <alias> <drive-id>",
		Short: "Create or update a drive alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Alias:   args[0],
				DriveID: args[1],
				Stdout:  cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("alias-set")
			return NewHandler(container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
