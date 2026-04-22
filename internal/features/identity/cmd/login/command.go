package login

import (
	"context"
	"github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateLoginCmd constructs and returns the cobra.Command for the login operation.
func CreateLoginCmd(container di.Container) *cobra.Command {
	opts := Options{}

	l, _ := container.Logger().CreateLogger("auth-login")
	handler := NewCommand(
		container.Config(),
		container.Identity(),
		l,
	)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		Long: `Authenticate with OneDrive using various methods (Interactive, Device Code, Service Principal).
You can specify the method via flags or in your profile configuration.`,
		Handler: handler,
		Options: &CommandContext{Options: opts},
		CtxFunc: func(ctx context.Context, c *CommandContext) *CommandContext {
			c.Ctx = ctx
			// Assuming stdout/stderr are stored in some way; this might need adjustments
			// if context doesn't provide them. For now, skipping direct I/O assignment.
			return c
		},
	})

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
