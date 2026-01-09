package di

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type ProfileService2 interface {
	PutProfile(context.Context, *azidentity.AuthenticationRecord) (string, error)
	GetProfile(context.Context, string) (*azidentity.AuthenticationRecord, error)
	DeleteProfile(context.Context, string) error
}
