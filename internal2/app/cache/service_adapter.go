package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/memory"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
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
	service2 domaincache.Service2
}

func memoryCacheFactory() abstractions.KeyValueStore {
	return memory.NewStore()
}

func BoltCacheFactory(path, bucket string) func() abstractions.KeyValueStore {
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

func NewServiceAdapter(authCachePath string, service2 domaincache.Service2) *ServiceAdapter {
	driveCacheStore := BoltCacheFactory(authCachePath, driveCacheName)()

	_ = service2.CreateCache(context.Background(), profileCacheName, memoryCacheFactory)
	_ = service2.CreateCache(context.Background(), configurationCacheName, memoryCacheFactory)
	_ = service2.CreateCache(context.Background(), driveCacheName, func() abstractions.KeyValueStore { return driveCacheStore })
	_ = service2.CreateCache(context.Background(), fileCacheName, siblingBoltFactory(driveCacheStore.(*bolt.Store), fileCacheName))
	_ = service2.CreateCache(context.Background(), authCacheName, siblingBoltFactory(driveCacheStore.(*bolt.Store), authCacheName))

	return &ServiceAdapter{
		service2: service2,
	}
}

func (s *ServiceAdapter) getProfileCache(ctx context.Context) (*abstractions.Cache2, error) {
	profileCache, exists := s.service2.GetCache(ctx, profileCacheName)
	if !exists {
		return nil, errors.New("No profile cache found")
	}

	return profileCache, nil
}

func (s *ServiceAdapter) getConfigurationCache(ctx context.Context) (*abstractions.Cache2, error) {
	configurationCache, exists := s.service2.GetCache(ctx, configurationCacheName)
	if !exists {
		return nil, errors.New("No configuration cache found")
	}

	return configurationCache, nil
}

func (s *ServiceAdapter) getDriveCache(ctx context.Context) (*abstractions.Cache2, error) {
	driveCache, exists := s.service2.GetCache(ctx, driveCacheName)
	if !exists {
		return nil, errors.New("No drive cache found")
	}

	return driveCache, nil
}

func (s *ServiceAdapter) getFileCache(ctx context.Context) (*abstractions.Cache2, error) {
	fileCache, exists := s.service2.GetCache(ctx, fileCacheName)
	if !exists {
		return nil, errors.New("No file cache found")
	}

	return fileCache, nil
}

func (s *ServiceAdapter) getAuthCache(ctx context.Context) (*abstractions.Cache2, error) {
	authCache, exists := s.service2.GetCache(ctx, authCacheName)
	if !exists {
		return nil, errors.New("No auth cache found")
	}

	return authCache, nil
}

// DeleteProfile implements [cache.CacheService].
func (s *ServiceAdapter) DeleteProfile(ctx context.Context, name string) error {
	authCache, err := s.getAuthCache(ctx)
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
	profileCache, err := s.getProfileCache(ctx)
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
	configurationCache, err := s.getConfigurationCache(ctx)
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
	driveCache, err := s.getDriveCache(ctx)
	if err != nil {
		return record, err
	}

	if err := driveCache.Get(ctx, func() ([]byte, error) { return json.Marshal(name) }, func(data []byte) error {
		parseNode, err := jsonserialization.NewJsonParseNode(data)
		if err != nil {
			return err
		}
		parsable, err := parseNode.GetObjectValue(domaincache.CreateCachedChildrenFromDiscriminatorValue)
		if err != nil {
			return err
		}
		cachedChildren, ok := parsable.(*cache.CachedChildren)
		if !ok {
			return errors.New("cached value is not of type CachedChildren")
		}
		record = *cachedChildren
		return nil
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return record, err
	}

	return record, nil
}

// GetItem implements [cache.CacheService].
func (s *ServiceAdapter) GetItem(ctx context.Context, name string) (cache.CachedItem, error) {
	var record cache.CachedItem
	fileCache, err := s.getFileCache(ctx)
	if err != nil {
		return record, err
	}

	if err := fileCache.Get(ctx, func() ([]byte, error) { return json.Marshal(name) }, func(data []byte) error {
		parseNode, err := jsonserialization.NewJsonParseNode(data)
		if err != nil {
			return err
		}
		parsable, err := parseNode.GetObjectValue(domaincache.CreateCachedItemFromDiscriminatorValue)
		if err != nil {
			return err
		}
		cachedItem, ok := parsable.(*cache.CachedItem)
		if !ok {
			return errors.New("cached value is not of type CachedItem")
		}
		record = *cachedItem
		return nil
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return record, err
	}

	return record, nil
}

// GetProfile implements [cache.CacheService].
func (s *ServiceAdapter) GetProfile(ctx context.Context, name string) (azidentity.AuthenticationRecord, error) {
	var record azidentity.AuthenticationRecord
	authCache, err := s.getAuthCache(ctx)
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
	profileCache, err := s.getProfileCache(ctx)
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
	configurationCache, err := s.getConfigurationCache(ctx)
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
	driveCache, err := s.getDriveCache(ctx)
	if err != nil {
		return err
	}

	if err := driveCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) {
		writer := jsonserialization.NewJsonSerializationWriter()

		if err := writer.WriteObjectValue("", &record); err != nil {
			return nil, err
		}

		return writer.GetSerializedContent()
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}

// SetItem implements [cache.CacheService].
func (s *ServiceAdapter) SetItem(ctx context.Context, name string, record cache.CachedItem) error {
	fileCache, err := s.getFileCache(ctx)
	if err != nil {
		return err
	}

	if err := fileCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) {
		writer := jsonserialization.NewJsonSerializationWriter()

		if err := writer.WriteObjectValue("", &record); err != nil {
			return nil, err
		}

		return writer.GetSerializedContent()
	}); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}

// SetProfile implements [cache.CacheService].
func (s *ServiceAdapter) SetProfile(ctx context.Context, name string, record azidentity.AuthenticationRecord) error {
	authCache, err := s.getAuthCache(ctx)
	if err != nil {
		return err
	}

	if err := authCache.Set(ctx, func() ([]byte, error) { return json.Marshal(name) }, func() ([]byte, error) { return json.Marshal(record) }); !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	return nil
}
