package current

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "current"
	loggerID    = "cli"
)

// CreateCurrentCmd constructs and returns the cobra.Command for the current operation.
func CreateCurrentCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewCurrentCmd(container).Run(cmd.Context(), opts)
		},
	}
}
