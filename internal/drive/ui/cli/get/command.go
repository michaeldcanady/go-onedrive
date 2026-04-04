package get

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateGetCmd constructs and returns the cobra.Command for the 'drive get' operation.
func CreateGetCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:               "get [drive-ref]",
		Short:             "Display details for a specific drive",
		Long:              "Provide a drive ID, name, or alias. If omitted, the currently active drive is shown.",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container, true),
		RunE: func(cmd *cobra.Command, args []string) error {
			driveRef := ""
			if len(args) > 0 {
				driveRef = args[0]
			}
			opts := Options{
				DriveRef: driveRef,
				Stdout:   cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("drive-get")
			return NewHandler(container.Drive(), container.Alias(), log).Handle(cmd.Context(), opts)
		},
	}
}
