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
	OfflineAccessScope     = "offline_access"

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
			logger, _ := container.LoggerService.GetLogger("cli")

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Load providers
			credProvider, err := container.CredentialProvider(ctx)
			if err != nil {
				logger.Error("failed to initialize credential provider", logging.String("error", err.Error()))
				return fmt.Errorf("failed to initialize credential provider: %w", err)
			}

			cacheService, err := container.CacheService(ctx)
			if err != nil {
				logger.Warn("cache service unavailable", logging.String("error", err.Error()))
			}

			profileName := container.Options.ProfileName

			// Load cached authentication record (if any)
			var record azidentity.AuthenticationRecord
			if cacheService != nil {
				rec, err := cacheService.GetProfile(ctx, profileName)
				if err == nil {
					record = rec
				} else {
					logger.Warn("no cached authentication record found", logging.String("error", err.Error()))
				}
			}

			// Create credential (may include cached record)
			cred, err := credProvider.Credential(ctx, profileName)
			if err != nil {
				logger.Error("failed to create credential", logging.String("error", err.Error()))
				return fmt.Errorf("failed to create credential: %w", err)
			}

			// Ensure credential supports explicit authentication
			authenticator, ok := cred.(interface {
				Authenticate(context.Context, *policy.TokenRequestOptions) (azidentity.AuthenticationRecord, error)
			})
			if !ok {
				logger.Error("configured credential does not support explicit authentication")
				return fmt.Errorf("configured credential does not support explicit authentication")
			}

			tokenOpts := &policy.TokenRequestOptions{
				Scopes: []string{
					FilesReadWriteAllScope,
					UserReadScope,
				},
				EnableCAE: true,
			}

			var token azcore.AccessToken
			var success bool

			for attempt := 0; attempt < maxAuthAttempts; attempt++ {
				needsAuth := force || isEmptyRecord(record)

				if needsAuth {
					logger.Info("Starting authentication flow...")

					newRecord, err := authenticator.Authenticate(ctx, tokenOpts)
					if err != nil {
						logger.Error("failed authentication", logging.String("error", err.Error()))
						return fmt.Errorf("authentication failed: %w", err)
					}

					record = newRecord

					// Save updated record
					if cacheService != nil {
						if err := cacheService.SetProfile(ctx, profileName, record); err != nil {
							logger.Warn("failed to cache authentication record", logging.String("error", err.Error()))
						}
					}

					// Reload credential with updated record
					cred, err = credProvider.Credential(ctx, profileName)
					if err != nil {
						logger.Error("failed to reload credential", logging.String("error", err.Error()))
						return fmt.Errorf("failed to reload credential: %w", err)
					}
				}

				token, err = cred.GetToken(ctx, *tokenOpts)
				if err != nil {
					if isAuthRequired(err) {
						logger.Warn("token request requires authentication; retrying")
						record = azidentity.AuthenticationRecord{} // force re-auth
						continue
					}
					logger.Error("failed to retrieve token", logging.String("error", err.Error()))
					return fmt.Errorf("failed to retrieve token: %w", err)
				}

				success = true
				break
			}

			if !success {
				logger.Error("authentication failed after max attempts", logging.Int("attempts", maxAuthAttempts))
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
