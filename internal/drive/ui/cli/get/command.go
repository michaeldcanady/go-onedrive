package get

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/spf13/cobra"
)

// CreateGetCmd constructs and returns the cobra.Command for the 'drive get' operation.
func CreateGetCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:               "get [drive-ref]",
		Short:             "Display details for a specific drive",
		Long:              "Provide a drive ID, name, or alias. If omitted, the currently active drive is shown.",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container, true),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.DriveRef = args[0]

			d, err := container.Drive().ResolveDrive(cmd.Context(), opts.DriveRef)
			if err != nil && !errors.Is(err, drive.ErrDriveNotFound) {
				return err
			}
			opts.DriveRef = d.ID
			if d == (drive.Drive{}) {
				if driveID, err := container.Alias().GetDriveIDByAlias(opts.DriveRef); err != nil {
					return coreerrors.NewNotFound(
						errors.New("drive ref not found"),
						fmt.Sprintf("unknown drive reference '%s'", opts.DriveRef),
						"Use 'odc drive list' to see available drives and aliases.",
					).WithContext(coreerrors.KeyName, opts.DriveRef)
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
