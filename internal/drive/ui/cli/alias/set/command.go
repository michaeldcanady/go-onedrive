package set

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateSetCmd constructs and returns the cobra.Command for the 'drive alias set' operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:               "set <drive-id> <alias>",
		Short:             "Create or update a drive alias",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: shared.ProviderPathCompletion(container, false),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate that the drive ID exists before attempting to set the alias
			driveID := args[0]
			_, err := container.Drive().ResolveDrive(cmd.Context(), driveID)
			if err != nil {
				return err
			}

			// validate that alias is not already in use for a different drive
			alias := args[1]
			existingDriveID, err := container.State().GetDriveAlias(alias)
			if err != nil {
				return err
			}
			if existingDriveID != "" && existingDriveID != driveID {
				return fmt.Errorf("alias '%s' is already in use for a different drive", alias)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Alias:   args[1],
				DriveID: args[0],
				Stdout:  cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("alias-set")
			return NewHandler(container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
