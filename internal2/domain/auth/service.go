package auth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AuthService interface {
	Login(ctx context.Context, profile string, opts LoginOptions) (*LoginResult, error)
	Logout(ctx context.Context, profile string, force bool) error

	azcore.TokenCredential
}
