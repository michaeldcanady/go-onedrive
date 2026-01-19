package auth

import (
	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth/login"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth/logout"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/cobra"
)

// CreateAuthCmd constructs the `auth` parent command, which groups all
// authentication-related subcommands (login, logout, token mgmt, etc).
func CreateAuthCmd(container *di.Container, logger logging.Logger) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for OneDrive",
		Long: `Authentication commands for managing your OneDrive session.

Use these commands to log in, log out, refresh tokens, or inspect your
current authentication state.
`,
	}

	// Subcommands
	authCmd.AddCommand(
		login.CreateLoginCmd(container, logger),
		logout.CreateLogoutCmd(),
	)

	return authCmd
}
