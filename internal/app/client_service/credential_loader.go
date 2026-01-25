package clientservice

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// credentialLoader loads a token credential
type credentialLoader interface {
	LoadCredential(ctx context.Context, name string) (azcore.TokenCredential, error)
}
