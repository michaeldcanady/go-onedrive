package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type Authenticator interface {
	Authenticate(ctx context.Context, opts *policy.TokenRequestOptions) (azidentity.AuthenticationRecord, error)
}

func InteractiveBrowserCredentialFactory(authConfig config.AuthenticationConfig) (*azidentity.InteractiveBrowserCredential, error) {
	return azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		TenantID: authConfig.GetTenantID(),
		ClientID: authConfig.GetClientID(),
		// Force users to use login command and authenticate that way.
		DisableAutomaticAuthentication: true,
		RedirectURL:                    authConfig.GetRedirectURI(),
	})
}

func CredentialFactory(authConfig config.AuthenticationConfig) (azcore.TokenCredential, error) {
	switch authConfig.GetAuthenticationType() {
	case "interactiveBrowser":
		return InteractiveBrowserCredentialFactory(authConfig)
	default:
		return nil, fmt.Errorf("Unsupported authentication type: %s", authConfig.GetAuthenticationType())
	}
}

func ClientFactory(config config.AuthenticationConfig) (*msgraphsdkgo.GraphServiceClient, error) {
	credential, err := CredentialFactory(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %w", err)
	}

	return msgraphsdkgo.NewGraphServiceClientWithCredentials(credential, []string{"Files.ReadWrite"})
}
