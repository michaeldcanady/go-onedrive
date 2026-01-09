package di

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	credentialservice "github.com/michaeldcanady/go-onedrive/internal/app/credential_service"
)

type CredentialService interface {
	LoadCredential(ctx context.Context, profile *azidentity.AuthenticationRecord) (azcore.TokenCredential, error)
	Authenticate(ctx context.Context, opts ...credentialservice.AuthenticationOption) (azcore.AccessToken, error)
	azcore.TokenCredential
}
