package clientservice

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// credentialLoader loads a token credential
type credentialLoader interface {
	LoadCredential(ctx context.Context, profile *azidentity.AuthenticationRecord) (azcore.TokenCredential, error)
	azcore.TokenCredential
}
