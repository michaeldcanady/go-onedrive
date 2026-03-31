package use

import (
	"github.com/michaeldcanady/go-onedrive/internal/feature/di"
	"github.com/spf13/cobra"
)

// CreateUseCmd constructs and returns the cobra.Command for the 'drive use' operation.
func CreateUseCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use <drive-ref>",
		Short: "Set the active OneDrive drive",
		Long:  "Specify a drive ID, name, or alias to be used for subsequent commands.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveRef: args[0],
				Stdout:   cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("drive-use")
			return NewHandler(container.Drive(), container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
