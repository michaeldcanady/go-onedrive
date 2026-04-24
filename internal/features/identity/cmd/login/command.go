package login

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateLoginCmd constructs and returns the cobra.Command for the login operation.
func CreateLoginCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("auth-login")
	handler := NewCommand(
		container.Config(),
		container.Identity(),
		l,
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		Long: `Authenticate with OneDrive using various methods (Interactive, Device Code, Service Principal).
You can specify the method via flags or in your profile configuration.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c = NewCommandContext(cmd.Context(), &opts)

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	cmd.Flags().StringVar(&opts.IdentityID, "id", "", "The specific identity (email) to authenticate")
	cmd.Flags().StringVar(&opts.Alias, "alias", "", "An optional human-friendly name for this identity")
	cmd.Flags().BoolVar(&opts.ShowToken, "show-token", false, "Display the access token after login")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force re-authentication even if a valid profile exists")
	cmd.Flags().StringVar(&opts.Method, "method", "", "Authentication method (interactive, device-code, client-secret, environment)")
	cmd.Flags().StringVar(&opts.TenantID, "tenant-id", "", "Azure AD tenant ID")
	cmd.Flags().StringVar(&opts.ClientID, "client-id", "", "Azure AD client ID")
	cmd.Flags().StringVar(&opts.ClientSecret, "client-secret", "", "Azure AD client secret (for client-secret method)")

	return cmd
}
