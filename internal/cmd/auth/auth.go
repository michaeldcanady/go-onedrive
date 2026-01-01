package auth

import (
	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth/login"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth/logout"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/cobra"
)

func CreateAuthCmd(logger logging.Logger, credentialService credentialService, profileService ProfileService) *cobra.Command {
	// authCmd represents the auth command
	var authCmd = &cobra.Command{
		Use:   "auth",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}

	authCmd.AddCommand(login.CreateLoginCmd(logger, credentialService, profileService))
	authCmd.AddCommand(logout.CreateLogoutCmd())

	return authCmd
}
