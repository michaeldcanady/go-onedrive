package app

import (
	"context"
	"path/filepath"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	appcache "github.com/michaeldcanady/go-onedrive/internal/cache/app"
	domaingraph "github.com/michaeldcanady/go-onedrive/internal/core/graph/domain"
	graphinfra "github.com/michaeldcanady/go-onedrive/internal/core/graph/infra"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	infrafile "github.com/michaeldcanady/go-onedrive/internal/fs/infra"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
	infrabolt "github.com/michaeldcanady/go-onedrive/pkg/cache/bolt"
)

func (c *Container) cacheStore() *infrabolt.Store {
	c.cacheStoreOnce.Do(func() {
		c.sharedStore = c.newCacheStore()
	})
	return c.sharedStore
}

func (c *Container) newCacheStore() *infrabolt.Store {
	environmentService := c.EnvironmentService()
	cachePath, _ := environmentService.CacheDir()
	dbPath := filepath.Join(cachePath, "cache.db")

	// Create a "root" store with a dummy bucket just to open the DB
	store, err := infrabolt.NewStore(dbPath, "root")
	if err != nil {
		panic(err) // Critical failure if we can't open the cache
	}
	return store
}

func (c *Container) clientProvider() domaingraph.ClientProvider {
	c.clientOnce.Do(func() {
		c.clientProvide = c.newClientProvider()
	})

	return c.clientProvide
}

func (c *Container) newClientProvider() domaingraph.ClientProvider {
	return graphinfra.New(c.Auth(), c.getLogger("graph"))
}

func (c *Container) metadataCache() infrafile.MetadataCache {
	c.metadataCacheOnce.Do(func() {
		c.metadataCacheCache = c.newMetadataCache()
	})
	return c.metadataCacheCache
}

func (c *Container) newMetadataCache() infrafile.MetadataCache {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "metadata", appcache.SiblingBoltFactory(c.cacheStore(), "metadata"))

	typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainfs.Metadata]{})
	return infrafile.NewMetadataCacheAdapter(typedCache)
}

func (c *Container) metadataListingCache() infrafile.ListingCache {
	c.metadataListingCacheOnce.Do(func() {
		c.metadataListingCacheCache = c.newMetadataListingCache()
	})
	return c.metadataListingCacheCache
}

func (c *Container) newMetadataListingCache() infrafile.ListingCache {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "metadatl", appcache.SiblingBoltFactory(c.cacheStore(), "metadatl"))

	typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainfs.Listing]{})
	return infrafile.NewMetadataListingCacheAdapter(typedCache)
}

func (c *Container) contentsCache() infrafile.ContentsCache {
	c.contentsCacheOnce.Do(func() {
		c.contentsCacheCache = c.newContentsCache()
	})
	return c.contentsCacheCache
}

func (c *Container) newContentsCache() infrafile.ContentsCache {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "contents", appcache.SiblingBoltFactory(c.cacheStore(), "contents"))

	typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainfs.Contents]{})
	return infrafile.NewContentsCacheAdapter(typedCache)
}

func (c *Container) pathIDCache() infrafile.PathIDCache {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "path_id", appcache.SiblingBoltFactory(c.cacheStore(), "path_id"))

	typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[string]{})
	return infrafile.NewPathIDCacheAdapter(typedCache)
}

func (c *Container) authCache() pkgcache.Cache[domainauth.AccessToken] {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "auth_tokens", appcache.SiblingBoltFactory(c.cacheStore(), "auth_tokens"))

	return appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainauth.AccessToken]{})
}

func (c *Container) accountCache() pkgcache.Cache[domainaccount.Account] {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "account", appcache.SiblingBoltFactory(c.cacheStore(), "account"))

	return appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainaccount.Account]{})
}
