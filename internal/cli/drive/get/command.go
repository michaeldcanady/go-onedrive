package get

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "get"
	loggerID    = "cli"
)

// CreateGetCmd constructs and returns the cobra.Command for the get operation.
func CreateGetCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id|alias>",
		Short: "Get information of named drive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveIDOrAlias: args[0],
				Stdout:         cmd.OutOrStdout(),
				Stderr:         cmd.ErrOrStderr(),
			}

			return NewGetCmd(container).Run(cmd.Context(), opts)
		},
	}
}
