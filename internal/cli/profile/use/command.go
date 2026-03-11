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
		Use:   "use <name>",
		Short: "Set current profile",
		Long: `You can set an existing profile as the current active profile for your
session. All subsequent commands will use the configuration and
authentication tokens associated with this profile.`,
		Example: `  # Switch to the 'personal' profile
  odc profile use personal`,
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Name:   args[0],
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewUseCmd(container).Run(cmd.Context(), opts)
		},
	}
}
