package login

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateLoginCmd constructs and returns the cobra.Command for the login operation.
func CreateLoginCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		Long: `Authenticate with OneDrive using various methods (Interactive, Device Code, Service Principal).
You can specify the method via flags or in your profile configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			l, _ := container.Logger().CreateLogger("auth-login")
			handler := NewHandler(
				container.Config(),
				container.Identity(),
				l,
			)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.ShowToken, "show-token", false, "Display the access token after login")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force re-authentication even if a valid profile exists")
	cmd.Flags().StringVar(&opts.Method, "method", "", "Authentication method (interactive, device-code, client-secret, environment)")
	cmd.Flags().StringVar(&opts.TenantID, "tenant-id", "", "Azure AD tenant ID")
	cmd.Flags().StringVar(&opts.ClientID, "client-id", "", "Azure AD client ID")
	cmd.Flags().StringVar(&opts.ClientSecret, "client-secret", "", "Azure AD client secret (for client-secret method)")

	return cmd
}
