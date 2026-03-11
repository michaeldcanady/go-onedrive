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
		Long: `You can authenticate with OneDrive to allow this application to access your
files. This command opens a browser window for you to log in to your
Microsoft account. Once authenticated, your access token is stored securely
in your active profile.`,
		Example: `  # Log in to your OneDrive account
  odc auth login

  # Force re-authentication even if you're already logged in
  odc auth login --force

  # Log in and display the access token in your terminal
  odc auth login --show-token`,
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
