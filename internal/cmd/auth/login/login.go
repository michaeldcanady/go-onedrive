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
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive using the configured authentication method",
		Long: `Authenticate with OneDrive using the authentication settings defined in your configuration.

This command initiates an authentication flow (browser, device code, etc.)
and stores the resulting authentication record for future CLI operations.`,
		Example: `
  # Perform interactive login
  go-onedrive auth login

  # Login and display the access token
  go-onedrive auth login --show-token
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			logger := container.Logger

			// Load existing profile
			logger.Info("Loading cached authentication profile...")
			profile, err := container.ProfileService.Load(ctx)
			if err != nil {
				logger.Warn("Unable to load profile", logging.String("error", err.Error()))
			}

			// Load credential provider
			cred, err := container.CredentialService.LoadCredential(ctx)
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
					OfflineAccessScope,
				},
				EnableCAE: true,
			}
			logger.Debug("Authentication options", logging.Any("options", *options))

			// Determine if authentication is required
			needsAuth := profile == nil || isEmptyRecord(*profile)
			if needsAuth {
				logger.Warn("No valid authentication profile found. Starting login flow...")

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

			// Optional: show token
			if showToken, _ := cmd.Flags().GetBool("show-token"); showToken {
				fmt.Printf("Access Token:\n%s\n", token.Token)
				return nil
			}

			// Save updated profile
			logger.Info("Saving authentication profile...")
			if err := container.ProfileService.Save(ctx, profile); err != nil {
				logger.Error("Unable to save profile", logging.String("error", err.Error()))
				return errors.Join(errors.New("unable to save profile"), err)
			}

			logger.Info("Login complete.")
			return nil
		},
	}

	cmd.Flags().Bool("show-token", false, "Display the access token after login")

	return cmd
}

// isEmptyRecord checks whether an AuthenticationRecord is effectively empty.
func isEmptyRecord(r azidentity.AuthenticationRecord) bool {
	return r.ClientID == "" &&
		r.TenantID == "" &&
		r.HomeAccountID == "" &&
		r.Username == ""
}
