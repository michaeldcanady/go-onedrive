package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/app"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/cobra"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
	SitesReadWriteAllScope = "Sites.ReadWrite.All"
	OfflineAccessScope     = "offline_access"
	AuthConfig             = "auth"
)

func CreateLoginCmd(container *di.Container) *cobra.Command {

	// loginCmd authenticates the user using the configured authentication method.
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive using the configured authentication method",
		Long: `Authenticate with OneDrive using the authentication settings defined in your configuration.

This command forces an authentication flow (e.g., interactive browser, device code)
and stores the resulting token for future CLI operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			container.Logger.Info("Loading cached profile...")
			profile, err := container.ProfileService.Load(ctx)
			if err != nil {
				container.Logger.Error("Unable to load profile", logging.String("error", err.Error()))
			}

			// Load credential from config (may use cached profile)
			cred, err := container.CredentialService.LoadCredential(ctx)
			if err != nil {
				container.Logger.Error("Failed to load credential from config", logging.String("error", err.Error()))
				return fmt.Errorf("failed to initialize credential: %w", err)
			}

			authenticator, ok := cred.(app.Authenticator)
			if !ok {
				container.Logger.Error("Configured credential does not support explicit authentication")
				return fmt.Errorf("configured credential does not support explicit authentication")
			}

			// Token request options
			options := &policy.TokenRequestOptions{
				Scopes: []string{
					FilesReadWriteAllScope,
					UserReadScope,
				},
				EnableCAE: true,
			}
			container.Logger.Debug("Authentication options", logging.Any("options", *options))

			// Determine if we need to authenticate
			needsAuth := profile == nil || *profile == (azidentity.AuthenticationRecord{})

			if needsAuth {
				container.Logger.Warn("No valid profile found. Starting authentication flow...")

				record, err := authenticator.Authenticate(ctx, options)
				if err != nil {
					container.Logger.Error("Authentication failed", logging.String("error", err.Error()))
					return fmt.Errorf("authentication failed: %w", err)
				}

				profile = &record
				container.Logger.Info("Authentication successful")
				container.Logger.Debug("Authentication record", logging.Any("profile", profile))
			}

			// Retrieve access token
			container.Logger.Info("Retrieving access token...")
			token, err := cred.GetToken(ctx, *options)
			if err != nil {
				container.Logger.Error("Failed to retrieve token", logging.String("error", err.Error()))
				return fmt.Errorf("failed to retrieve token: %w", err)
			}
			container.Logger.Info("Access token retrieved successfully")
			container.Logger.Debug("Access token", logging.Any("token", token))

			// Optional flag to show token
			if showToken, _ := cmd.Flags().GetBool("show-token"); showToken {
				fmt.Printf("Access Token: %s\n", token.Token)
				return nil
			}

			// Save updated profile
			container.Logger.Info("Saving authentication profile...")
			if err := container.ProfileService.Save(ctx, profile); err != nil {
				container.Logger.Error("Unable to save profile", logging.String("error", err.Error()))
				return errors.Join(errors.New("unable to save profile"), err)
			}

			container.Logger.Info("Login complete.")
			return nil
		},
	}

	// Optional flag to show token (safer default)
	loginCmd.Flags().Bool("show-token", false, "Display the access token after login")

	return loginCmd
}
