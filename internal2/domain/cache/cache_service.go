package cache

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type CacheService interface {
	// GetProfile returns the currently cached profile by name.
	GetProfile(ctx context.Context, name string) (azidentity.AuthenticationRecord, error)
	// SetProfile caches the provided profile by name.
	SetProfile(ctx context.Context, name string, record azidentity.AuthenticationRecord) error
	DeleteProfile(ctx context.Context, name string) error

	GetConfiguration(ctx context.Context, name string) (config.Configuration3, error)
	SetConfiguration(ctx context.Context, name string, record config.Configuration3) error

	GetCLIProfile(ctx context.Context, name string) (profile.Profile, error)
	SetCLIProfile(ctx context.Context, name string, profile profile.Profile) error

	GetDrive(ctx context.Context, name string) (CachedChildren, error)
	SetDrive(ctx context.Context, name string, record CachedChildren) error

	GetItem(ctx context.Context, name string) (CachedItem, error)
	SetItem(ctx context.Context, name string, record CachedItem) error
}
