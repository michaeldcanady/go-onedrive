package login

import (
	"context"
	"fmt"

	credentialservice "github.com/michaeldcanady/go-onedrive/internal/app/credential_service"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/cobra"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
	SitesReadWriteAllScope = "Sites.ReadWrite.All"
	OfflineAccessScope     = "offline_access"

	showTokenLongFlag = "show-token"
	showTokenUsage    = "Display the access token after login"

	forceLongFlag  = "force"
	forceShortFlag = "f"
	forceUsage     = "Force re-authentication even if a valid profile exists"
)

func CreateLoginCmd(container *di.Container) *cobra.Command {
	var (
		showToken bool
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		Long:  "Authenticate with OneDrive using the Microsoft identity platform.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			authenticationOpts := []credentialservice.AuthenticationOption{
				credentialservice.WithScopes(FilesReadWriteAllScope, UserReadScope),
				credentialservice.WithCAE(),
			}

			if force {
				authenticationOpts = append(authenticationOpts, credentialservice.WithForceAuthentication())
			}

			token, err := container.CredentialService.Authenticate(
				ctx,
				authenticationOpts...,
			)
			if err != nil {
				container.Logger.Error("authentication failed", logging.Any("error", err))
				return fmt.Errorf("authentication failed: %w", err)
			}

			if showToken {
				fmt.Printf("Access Token:\n%s\n", token.Token)
			}

			container.Logger.Info("Login complete.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&showToken, showTokenLongFlag, false, showTokenUsage)
	cmd.Flags().BoolVarP(&force, forceLongFlag, forceShortFlag, false, forceUsage)
	return cmd
}
