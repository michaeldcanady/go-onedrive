package logout

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "logout"
	loggerID    = "cli"

	forceLongFlag  = "force"
	forceShortFlag = "f"
	forceUsage     = "Force logout even if no active session is detected"
)

// CreateLogoutCmd constructs and returns the cobra.Command for the logout operation.
func CreateLogoutCmd(container didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of the current OneDrive profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			return NewLogoutCmd(container).Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}
