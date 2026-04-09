package set

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	aliaspkg "github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
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
			opts.DriveID = args[0]
			if d, err := container.Drive().ResolveDrive(cmd.Context(), opts.DriveID); err != nil && !errors.Is(err, drive.ErrDriveNotFound) {
				return err
			} else if d == (drive.Drive{}) {
				return coreerrors.NewNotFound(
					errors.New("drive not found"),
					"drive not found",
					"Use 'odc drive list' to see available drives.",
				).WithContext(coreerrors.KeyName, opts.DriveID)
			} else {
				opts.DriveID = d.ID
			}

			// validate that alias is not already in use for a different drive
			opts.Alias = args[1]
			existingDriveID, err := container.Alias().GetDriveIDByAlias(opts.Alias)
			if err != nil && !errors.Is(err, aliaspkg.ErrDriveIDNotFound) {
				return err
			}
			if existingDriveID != "" && existingDriveID != opts.DriveID {
				return coreerrors.NewConflict(
					errors.New("alias already in use"),
					fmt.Sprintf("alias '%s' is already in use for drive ID '%s'", opts.Alias, existingDriveID),
					"Please choose a different alias or remove the existing one.",
				).WithContext(coreerrors.KeyDriveID, existingDriveID)
			}

			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.OutOrStderr()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
