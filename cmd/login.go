package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/app"
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

var profileService app.ProfileService

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

		logger.Info("Loading cached profile...")
		profile, err := profileService.Load(ctx)
		if err != nil {
			logger.Error("Unable to load profile", logging.String("error", err.Error()))
		}

		// Load credential from config (may use cached profile)
		cred, err := credentialService.LoadCredential(ctx)
		if err != nil {
			logger.Error("Failed to load credential from config", logging.String("error", err.Error()))
			return fmt.Errorf("failed to initialize credential: %w", err)
		}

		authenticator, ok := cred.(app.Authenticator)
		if !ok {
			logger.Error("Configured credential does not support explicit authentication")
			return fmt.Errorf("configured credential does not support explicit authentication")
		}

		// Token request options
		options := &policy.TokenRequestOptions{
			Scopes: []string{
				FilesReadWriteAllScope,
				UserReadScope,
			},
		}
		logger.Debug("Authentication options", logging.Any("options", *options))

		// Determine if we need to authenticate
		needsAuth := profile == nil || *profile == (azidentity.AuthenticationRecord{})

		if needsAuth {
			logger.Warn("No valid profile found. Starting authentication flow...")

			record, err := authenticator.Authenticate(ctx, options)
			if err != nil {
				logger.Error("Authentication failed", logging.String("error", err.Error()))
				return fmt.Errorf("authentication failed: %w", err)
			}

			profile = &record
			logger.Info("Authentication successful")
			logger.Debug("Authentication record", logging.Any("profile", profile))
		}

		// Retrieve access token
		logger.Info("Retrieving access token...")
		token, err := cred.GetToken(ctx, *options)
		if err != nil {
			logger.Error("Failed to retrieve token", logging.String("error", err.Error()))
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		logger.Info("Access token retrieved successfully")
		logger.Debug("Access token", logging.Any("token", token))

		// Optional flag to show token
		if showToken, _ := cmd.Flags().GetBool("show-token"); showToken {
			fmt.Printf("Access Token: %s\n", token.Token)
			return nil
		}

		// Save updated profile
		logger.Info("Saving authentication profile...")
		if err := profileService.Save(ctx, profile); err != nil {
			logger.Error("Unable to save profile", logging.String("error", err.Error()))
			return errors.Join(errors.New("unable to save profile"), err)
		}

		logger.Info("Login complete.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	// Optional flag to show token (safer default)
	loginCmd.Flags().Bool("show-token", false, "Display the access token after login")
}
