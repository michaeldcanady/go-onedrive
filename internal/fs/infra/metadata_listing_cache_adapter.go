package infra

import (
	"context"

	domaincache "github.com/michaeldcanady/go-onedrive/internal/cache/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

var _ ListingCache = (*MetadataListingCacheAdapter)(nil)

type MetadataListingCacheAdapter struct {
	cache domaincache.Cache[domainfs.Listing]
}

func NewMetadataListingCacheAdapter(cache domaincache.Cache[domainfs.Listing]) *MetadataListingCacheAdapter {
	return &MetadataListingCacheAdapter{cache: cache}
}

func (a *MetadataListingCacheAdapter) Get(ctx context.Context, path string) (*domainfs.Listing, bool) {
	listing, err := a.cache.Get(ctx, path)
	if err != nil {
		return nil, false
	}
	return &listing, true
}

func (a *MetadataListingCacheAdapter) Put(ctx context.Context, path string, listing *domainfs.Listing) error {
	return a.cache.Set(ctx, path, *listing)
}

func (a *MetadataListingCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return a.cache.Delete(ctx, path)
}
