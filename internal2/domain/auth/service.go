package auth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AuthService interface {
	Login(ctx context.Context, profileName string, opts LoginOptions) (*LoginResult, error)
	//Logout(ctx context.Context, profileName string) error

	azcore.TokenCredential
}
