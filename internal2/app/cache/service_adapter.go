package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/memory"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

var _ cache.CacheService = (*ServiceAdapter)(nil)

const (
	authCacheName          = "auth"
	profileCacheName       = "profile"
	configurationCacheName = "configuration"
	driveCacheName         = "drive"
	fileCacheName          = "file"
)

type ServiceAdapter struct {
	service2 *Service2
}

func memoryCacheFactory() abstractions.KeyValueStore {
	return memory.NewStore()
}

func boltCacheFactory(path, bucket string) func() abstractions.KeyValueStore {
	return func() abstractions.KeyValueStore {
		store, err := bolt.NewStore(path, bucket)
		if err != nil {
			return nil
		}
		return store
	}
}

func siblingBoltFactory(store *bolt.Store, bucket string) func() abstractions.KeyValueStore {
	return func() abstractions.KeyValueStore {
		siblingStore, err := bolt.NewSiblingStore(store, bucket)
		if err != nil {
			return nil
		}
		return siblingStore
	}
}

func NewServiceAdapter(authCachePath, driveCachePath, fileCachePath string, service2 *Service2) *ServiceAdapter {
	driveCacheStore := boltCacheFactory(driveCachePath, driveCacheName)()

	_ = service2.CreateCache(profileCacheName, memoryCacheFactory)
	_ = service2.CreateCache(configurationCacheName, memoryCacheFactory)
	_ = service2.CreateCache(driveCacheName, func() abstractions.KeyValueStore { return driveCacheStore })
	_ = service2.CreateCache(fileCacheName, siblingBoltFactory(driveCacheStore.(*bolt.Store), fileCacheName))
	_ = service2.CreateCache(authCacheName, siblingBoltFactory(driveCacheStore.(*bolt.Store), authCacheName))

	return &ServiceAdapter{
		service2: service2,
	}
}

func (s *ServiceAdapter) getProfileCache() (*abstractions.Cache2, error) {
	profileCache, exists := s.service2.GetCache(profileCacheName)
	if !exists {
		return nil, errors.New("No profile cache found")
	}

	return profileCache, nil
}

func (s *ServiceAdapter) getConfigurationCache() (*abstractions.Cache2, error) {
	configurationCache, exists := s.service2.GetCache(configurationCacheName)
	if !exists {
		return nil, errors.New("No configuration cache found")
	}

	return configurationCache, nil
}

func (s *ServiceAdapter) getDriveCache() (*abstractions.Cache2, error) {
	driveCache, exists := s.service2.GetCache(driveCacheName)
	if !exists {
		return nil, errors.New("No drive cache found")
	}

	return driveCache, nil
}

func (s *ServiceAdapter) getFileCache() (*abstractions.Cache2, error) {
	fileCache, exists := s.service2.GetCache(fileCacheName)
	if !exists {
		return nil, errors.New("No file cache found")
	}

	return fileCache, nil
}

func (s *ServiceAdapter) getAuthCache() (*abstractions.Cache2, error) {
	authCache, exists := s.service2.GetCache(authCacheName)
	if !exists {
		return nil, errors.New("No auth cache found")
	}

	return authCache, nil
}

// DeleteProfile implements [cache.CacheService].
func (s *ServiceAdapter) DeleteProfile(ctx context.Context, name string) error {
	authCache, err := s.getAuthCache()
	if err != nil {
		return err
	}

	if err := authCache.Delete(ctx, func() ([]byte, error) { return json.Marshal(name) }); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}

	return nil
}

// GetCLIProfile implements [cache.CacheService].
func (s *ServiceAdapter) GetCLIProfile(ctx context.Context, name string) (profile.Profile, error) {
	var profile profile.Profile
	profileCache, err := s.getProfileCache()
	if err != nil {
		return profile, err
	}

	if err := profileCache.Get(ctx, func() ([]byte, error) { return json.Marshal(name) }, func(data []byte) error {
		if err := json.Unmarshal(data, &profile); err != nil {
			return err
		}
		return nil
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return profile, err
	}

	return profile, nil
}

// GetConfiguration implements [cache.CacheService].
func (s *ServiceAdapter) GetConfiguration(ctx context.Context, name string) (config.Configuration3, error) {
	var record config.Configuration3
	configurationCache, err := s.getConfigurationCache()
	if err != nil {
		return record, err
	}

	if err := configurationCache.Get(ctx, func() ([]byte, error) { return json.Marshal(name) }, func(data []byte) error {
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
		return nil
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return record, err
	}

	return record, nil
}

// GetDrive implements [cache.CacheService].
func (s *ServiceAdapter) GetDrive(ctx context.Context, name string) (cache.CachedChildren, error) {
	var record cache.CachedChildren
	driveCache, err := s.getDriveCache()
	if err != nil {
		return record, err
	}

	if err := driveCache.Get(ctx, func() ([]byte, error) { return json.Marshal(name) }, func(data []byte) error {
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
		return nil
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return record, err
	}

	return record, nil
}

// GetItem implements [cache.CacheService].
func (s *ServiceAdapter) GetItem(ctx context.Context, name string) (cache.CachedItem, error) {
	var record cache.CachedItem
	fileCache, err := s.getFileCache()
	if err != nil {
		return record, err
	}

	if err := fileCache.Get(ctx, func() ([]byte, error) { return json.Marshal(name) }, func(data []byte) error {
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
		return nil
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return record, err
	}

	return record, nil
}

// GetProfile implements [cache.CacheService].
func (s *ServiceAdapter) GetProfile(ctx context.Context, name string) (azidentity.AuthenticationRecord, error) {
	var record azidentity.AuthenticationRecord
	authCache, err := s.getAuthCache()
	if err != nil {
		return record, err
	}

	if err := authCache.Get(ctx, func() ([]byte, error) { return json.Marshal(name) }, func(data []byte) error {
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
		return nil
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return record, err
	}

	return record, nil
}

// SetCLIProfile implements [cache.CacheService].
func (s *ServiceAdapter) SetCLIProfile(ctx context.Context, name string, profile profile.Profile) error {
	profileCache, err := s.getProfileCache()
	if err != nil {
		return err
	}

	if err := profileCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) { return json.Marshal(profile) }); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}

// SetConfiguration implements [cache.CacheService].
func (s *ServiceAdapter) SetConfiguration(ctx context.Context, name string, record config.Configuration3) error {
	configurationCache, err := s.getConfigurationCache()
	if err != nil {
		return err
	}

	if err := configurationCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) { return json.Marshal(record) }); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}

// SetDrive implements [cache.CacheService].
func (s *ServiceAdapter) SetDrive(ctx context.Context, name string, record cache.CachedChildren) error {
	driveCache, err := s.getDriveCache()
	if err != nil {
		return err
	}

	if err := driveCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) { return json.Marshal(record) }); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}

// SetItem implements [cache.CacheService].
func (s *ServiceAdapter) SetItem(ctx context.Context, name string, record cache.CachedItem) error {
	fileCache, err := s.getFileCache()
	if err != nil {
		return err
	}

	if err := fileCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) { return json.Marshal(record) }); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}

// SetProfile implements [cache.CacheService].
func (s *ServiceAdapter) SetProfile(ctx context.Context, name string, record azidentity.AuthenticationRecord) error {
	authCache, err := s.getAuthCache()
	if err != nil {
		return err
	}

	if err := authCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) { return json.Marshal(record) }); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}
