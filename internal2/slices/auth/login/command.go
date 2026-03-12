package login

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// CreateLoginCmd constructs and returns the cobra.Command for the login operation.
func CreateLoginCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		Long: `You can authenticate with OneDrive to allow this application to access your
files. This command opens a browser window for you to log in to your
Microsoft account. Once authenticated, your access token is stored securely
in your active profile.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			l, _ := container.Logger().CreateLogger("auth-login")
			handler := NewHandler(
				container.Config(),
				container.State(),
				container.Identity(),
				l,
			)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.ShowToken, "show-token", false, "Display the access token after login")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force re-authentication even if a valid profile exists")

	return cmd
}
