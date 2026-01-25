package di

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type CredentialService interface {
	LoadCredential(ctx context.Context, profile string) (azcore.TokenCredential, error)
}
