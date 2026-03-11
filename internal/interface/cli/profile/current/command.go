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
		Long: `You can display the name of the profile that's currently active in your
session. This is useful for confirming which account you're interacting with.`,
		Example: `  # Show the name of the active profile
  odc profile current`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewCurrentCmd(container).Run(cmd.Context(), opts)
		},
	}
}
