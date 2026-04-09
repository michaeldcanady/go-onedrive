package set

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	aliaspkg "github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateSetCmd constructs and returns the cobra.Command for the 'drive alias set' operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:               "set <drive-id> <alias>",
		Short:             "Create or update a drive alias",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: shared.ProviderPathCompletion(container, false),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate that the drive ID exists before attempting to set the alias
			driveID := args[0]
			//_, err := container.Drive().ResolveDrive(cmd.Context(), driveID)
			//if err != nil {
			//	return err
			//}

			// validate that alias is not already in use for a different drive
			alias := args[1]
			existingDriveID, err := container.Alias().GetDriveIDByAlias(alias)
			if err != nil && !errors.Is(err, aliaspkg.ErrDriveIDNotFound) {
				return err
			}
			if existingDriveID != "" && existingDriveID != driveID {
				return fmt.Errorf("alias '%s' is already in use for a different drive", alias)
			}

			opts.Alias = alias
			opts.DriveID = driveID
			opts.Stdout = cmd.OutOrStdout()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
