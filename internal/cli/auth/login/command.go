package login

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	showTokenLongFlag = "show-token"
	showTokenUsage    = "Display the access token after login"

	forceLongFlag  = "force"
	forceShortFlag = "f"
	forceUsage     = "Force re-authentication even if a valid profile exists"

	commandName = "login"
	loggerID    = "cli"
)

// CreateLoginCmd constructs and returns the cobra.Command for the login operation.
func CreateLoginCmd(container didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			return NewLoginCmd(container).Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.ShowToken, showTokenLongFlag, false, showTokenUsage)
	cmd.Flags().BoolVarP(&opts.Force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}
