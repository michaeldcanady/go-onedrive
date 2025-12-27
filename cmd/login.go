package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
	SitesReadWriteAllScope = "Sites.ReadWrite.All"
	OfflineAccessScope     = "offline_access"
	AuthConfig             = "auth"
)

// loginCmd authenticates the user using the configured authentication method.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with OneDrive using the configured authentication method",
	Long: `Authenticate with OneDrive using the authentication settings defined in your configuration.

This command forces an authentication flow (e.g., interactive browser, device code)
and stores the resulting token for future CLI operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cred, err := loadCredentialFromConfig()
		if err != nil {
			logger.Error("Failed to load credential from config", logging.String("error", err.Error()))
			return fmt.Errorf("failed to initialize credential: %w", err)
		}

		authenticator, ok := cred.(Authenticator)
		if !ok {
			logger.Error("Configured credential does not support explicit authentication")
			return fmt.Errorf("configured credential does not support explicit authentication")
		}

		logger.Info("Starting authentication...")

		options := &policy.TokenRequestOptions{
			Scopes: []string{
				FilesReadWriteAllScope,
				UserReadScope,
			},
		}
		logger.Debug("authentication options", logging.Any("options", *options))

		logger.Info("Sending authentication request...")
		record, err := authenticator.Authenticate(
			context.Background(),
			options,
		)
		logger.Debug("authentication record", logging.Any("record", record))

		if err != nil {
			logger.Error("Authentication failed", logging.String("error", err.Error()))
			return fmt.Errorf("authentication failed: %w", err)
		}
		logger.Info("authentication successful")

		logger.Info("Retrieving access token...")
		token, err := cred.GetToken(
			context.Background(),
			*options,
		)
		logger.Debug("access token", logging.Any("token", token))

		if err != nil {
			logger.Error("Failed to retrieve token", logging.String("error", err.Error()))
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		logger.Info("Access token retrieved successfully")

		if showToken, _ := cmd.Flags().GetBool("show-token"); showToken {
			fmt.Printf("Access Token: %s\n", token.Token)
			return nil
		}

		// TODO: Securely store the token for future use.

		logger.Info("Login complete.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	// Optional flag to show token (safer default)
	loginCmd.Flags().Bool("show-token", false, "Display the access token after login")
}

// loadCredentialFromConfig reads the auth config and constructs the appropriate credential.
func loadCredentialFromConfig() (azcore.TokenCredential, error) {
	var authCfg config.AuthenticationConfigImpl

	sub := viper.Sub(AuthConfig)
	if sub == nil {
		return nil, fmt.Errorf("missing '%s' section in configuration", AuthConfig)
	}

	if err := sub.Unmarshal(&authCfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", AuthConfig, err)
	}

	cred, err := CredentialFactory(&authCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	return cred, nil
}
