package auth

import (
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/auth/login"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/auth/logout"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

// CreateAuthCmd constructs the `auth` parent command, which groups all
// authentication-related subcommands (login, logout, token mgmt, etc).
func CreateAuthCmd(container didomain.Container) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for OneDrive",
		Long: `Authentication commands for managing your OneDrive session.

Use these commands to log in, log out, refresh tokens, or inspect your
current authentication domainstate.
`,
	}

	// Subcommands
	authCmd.AddCommand(
		login.CreateLoginCmd(container),
		logout.CreateLogoutCmd(container),
	)

	return authCmd
}
