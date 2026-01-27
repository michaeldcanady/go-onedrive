package graph

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type CredentialProvider interface {
	Credential(ctx context.Context, profile string) (azcore.TokenCredential, error)
}
