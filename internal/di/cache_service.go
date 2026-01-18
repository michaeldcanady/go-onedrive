package di

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type CacheService interface {
	GetProfile(context.Context, string) (azidentity.AuthenticationRecord, error)
	SetProfile(ctx context.Context, name string, record azidentity.AuthenticationRecord) error
}
