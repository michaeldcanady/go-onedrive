package login

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/di2"
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

func CreateLoginCmd(container *di2.Container) *cobra.Command {
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

			logger := container.Logger().Must("cli")
			credProvider := container.CredentialProvider()
			cache := container.Cache()
			profileSvc := container.Profile()
			profileName := container.Options().ProfileName

			// Load cached record
			var record azidentity.AuthenticationRecord
			if cache != nil {
				rec, err := cache.GetProfile(ctx, profileName)
				if err == nil {
					record = rec
				}
			}

			// Create credential
			cred, err := credProvider.Credential(ctx, profileName)
			if err != nil {
				return fmt.Errorf("failed to create credential: %w", err)
			}

			authenticator, ok := cred.(interface {
				Authenticate(context.Context, *policy.TokenRequestOptions) (azidentity.AuthenticationRecord, error)
			})
			if !ok {
				return fmt.Errorf("credential does not support explicit authentication")
			}

			tokenOpts := &policy.TokenRequestOptions{
				Scopes:    []string{FilesReadWriteAllScope, UserReadScope},
				EnableCAE: true,
			}

			var token azcore.AccessToken
			for attempt := 0; attempt < maxAuthAttempts; attempt++ {
				needsAuth := force || isEmptyRecord(record)

				if needsAuth {
					newRecord, err := authenticator.Authenticate(ctx, tokenOpts)
					if err != nil {
						return fmt.Errorf("authentication failed: %w", err)
					}

					record = newRecord
					cache.SetProfile(ctx, profileName, record)

					cred, err = credProvider.Credential(ctx, profileName)
					if err != nil {
						return fmt.Errorf("failed to reload credential: %w", err)
					}
				}

				token, err = cred.GetToken(ctx, *tokenOpts)
				if err != nil {
					if isAuthRequired(err) {
						record = azidentity.AuthenticationRecord{}
						continue
					}
					return fmt.Errorf("failed to retrieve token: %w", err)
				}

				break
			}

			if showToken {
				fmt.Println("Access Token:")
				fmt.Println(token.Token)
			}

			logger.Info("Login complete.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&showToken, showTokenLongFlag, false, showTokenUsage)
	cmd.Flags().BoolVarP(&force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}
