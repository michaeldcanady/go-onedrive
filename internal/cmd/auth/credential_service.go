package auth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type credentialService interface {
	LoadCredential(ctx context.Context) (azcore.TokenCredential, error)
}
