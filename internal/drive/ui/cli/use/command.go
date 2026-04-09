package use

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/spf13/cobra"
)

// CreateUseCmd constructs and returns the cobra.Command for the 'drive use' operation.
func CreateUseCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:               "use <drive-ref>",
		Short:             "Set the active OneDrive drive",
		Long:              "Specify a drive ID, name, or alias to be used for subsequent commands.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container, true),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			aliasOrRef := args[0]

			d, err := container.Drive().ResolveDrive(cmd.Context(), aliasOrRef)
			if err != nil && !errors.Is(err, drive.ErrDriveNotFound) {
				return err
			}
			opts.DriveRef = d.ID
			if d == (drive.Drive{}) {
				if driveID, err := container.Alias().GetDriveIDByAlias(aliasOrRef); err != nil {
					return coreerrors.NewNotFound(
						errors.New("drive ref not found"),
						fmt.Sprintf("unknown drive reference '%s'", aliasOrRef),
						"Use 'odc drive list' to see available drives and aliases.",
					).WithContext(coreerrors.KeyName, aliasOrRef)
				} else {
					opts.DriveRef = driveID
				}
			}

			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Drive(), container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
