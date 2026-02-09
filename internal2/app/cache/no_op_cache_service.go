package cache

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

var _ (cache.CacheService) = (*NoopCacheService)(nil)

type NoopCacheService struct{}

func NewNoopCacheService() *NoopCacheService {
	return &NoopCacheService{}
}

// DeleteProfile implements [cache.CacheService].
func (n *NoopCacheService) DeleteProfile(_ context.Context, _ string) error {
	return cache.ErrUnavailableCache
}

// GetCLIProfile implements [cache.CacheService].
func (n *NoopCacheService) GetCLIProfile(_ context.Context, _ string) (profile.Profile, error) {
	return profile.Profile{}, cache.ErrUnavailableCache
}

// GetConfiguration implements [cache.CacheService].
func (n *NoopCacheService) GetConfiguration(_ context.Context, _ string) (config.Configuration3, error) {
	return config.Configuration3{}, cache.ErrUnavailableCache
}

// GetDrive implements [cache.CacheService].
func (n *NoopCacheService) GetDrive(_ context.Context, _ string) (cache.CachedChildren, error) {
	return cache.CachedChildren{}, cache.ErrUnavailableCache
}

// GetItem implements [cache.CacheService].
func (n *NoopCacheService) GetItem(_ context.Context, _ string) (cache.CachedItem, error) {
	return cache.CachedItem{}, cache.ErrUnavailableCache
}

// GetProfile implements [cache.CacheService].
func (n *NoopCacheService) GetProfile(_ context.Context, _ string) (azidentity.AuthenticationRecord, error) {
	return azidentity.AuthenticationRecord{}, cache.ErrUnavailableCache
}

// SetCLIProfile implements [cache.CacheService].
func (n *NoopCacheService) SetCLIProfile(_ context.Context, _ string, _ profile.Profile) error {
	return cache.ErrUnavailableCache
}

// SetConfiguration implements [cache.CacheService].
func (n *NoopCacheService) SetConfiguration(_ context.Context, _ string, _ config.Configuration3) error {
	return cache.ErrUnavailableCache
}

// SetDrive implements [cache.CacheService].
func (n *NoopCacheService) SetDrive(_ context.Context, _ string, _ cache.CachedChildren) error {
	return cache.ErrUnavailableCache
}

// SetItem implements [cache.CacheService].
func (n *NoopCacheService) SetItem(_ context.Context, _ string, _ cache.CachedItem) error {
	return cache.ErrUnavailableCache
}

// SetProfile implements [cache.CacheService].
func (n *NoopCacheService) SetProfile(_ context.Context, _ string, _ azidentity.AuthenticationRecord) error {
	return cache.ErrUnavailableCache
}
