package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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

	maxAuthAttempts = 3

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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			logger := container.Logger

			// Load existing profile
			record, err := container.CacheService.GetProfile(ctx, "default")
			if err != nil {
				logger.Warn("Unable to load profile", logging.String("error", err.Error()))
			}
			profile := &record

			// Load credential provider (reactive chain starts here)
			cred, err := container.CredentialService.LoadCredential(ctx)
			if err != nil {
				return fmt.Errorf("failed to initialize credential: %w", err)
			}

			authenticator, ok := cred.(Authenticator)
			if !ok {
				return fmt.Errorf("configured credential does not support explicit authentication")
			}

			options := &policy.TokenRequestOptions{
				Scopes: []string{
					FilesReadWriteAllScope,
					UserReadScope,
					//OfflineAccessScope,
				},
				EnableCAE: true,
			}

			var token azcore.AccessToken
			var success bool

			for range maxAuthAttempts {
				needsAuth := (profile == nil || isEmptyRecord(*profile)) || force
				if needsAuth {
					logger.Info("Starting authentication flow...")

					record, err := authenticator.Authenticate(ctx, options)
					if err != nil {
						logger.Error("authentication failed", logging.String("error", err.Error()))
						return fmt.Errorf("authentication failed: %w", err)
					}

					if err := container.CacheService.SetProfile(ctx, "default", record); err != nil {
						return fmt.Errorf("unable to save profile: %w", err)
					}

					profile = &record

					cred, err = container.CredentialService.LoadCredential(ctx)
					if err != nil {
						return fmt.Errorf("failed to reload credential: %w", err)
					}
				}

				token, err = cred.GetToken(ctx, *options)
				if err != nil {
					if isAuthRequired(err) {
						logger.Warn("Authentication required to obtain access token; clearing cached profile",
							logging.String("error", err.Error()))
						profile = nil
						continue
					}
					return fmt.Errorf("failed to retrieve token: %w", err)
				}

				success = true
				break
			}

			if !success {
				return fmt.Errorf("authentication failed after %d attempts", maxAuthAttempts)
			}

			if showToken {
				fmt.Printf("Access Token:\n%s\n", token.Token)
			}

			logger.Info("Login complete.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&showToken, showTokenLongFlag, false, showTokenUsage)
	cmd.Flags().BoolVarP(&force, forceLongFlag, forceShortFlag, false, forceUsage)
	return cmd
}

// isEmptyRecord checks whether an AuthenticationRecord is effectively empty.
func isEmptyRecord(r azidentity.AuthenticationRecord) bool {
	return r.ClientID == "" &&
		r.TenantID == "" &&
		r.HomeAccountID == "" &&
		r.Username == ""
}

func isAuthRequired(err error) bool {
	var authErr *azidentity.AuthenticationRequiredError
	return errors.As(err, &authErr)
}
