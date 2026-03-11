package use

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "use"
	loggerID    = "cli"
)

// CreateUseCmd constructs and returns the cobra.Command for the use operation.
func CreateUseCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use [drive-id|alias]",
		Short: "Sets the active drive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveIDOrAlias: args[0],
				Stdout:         cmd.OutOrStdout(),
				Stderr:         cmd.ErrOrStderr(),
			}

			return NewUseCmd(container).Run(cmd.Context(), opts)
		},
	}
}
