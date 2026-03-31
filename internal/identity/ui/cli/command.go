// Package auth provides the authentication-related CLI commands.
package auth

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/identity/ui/cli/login"
	"github.com/michaeldcanady/go-onedrive/internal/identity/ui/cli/logout"
	"github.com/spf13/cobra"
)

// CreateAuthCmd constructs and returns the cobra.Command for the 'auth' parent command.
func CreateAuthCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication with OneDrive",
		Long:  `Authenticate with OneDrive and manage your login state.`,
	}

	cmd.AddCommand(
		login.CreateLoginCmd(container),
		logout.CreateLogoutCmd(container),
	)

	return cmd
}
