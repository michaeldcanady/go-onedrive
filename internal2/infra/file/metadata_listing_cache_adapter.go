package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

var _ ListingCache = (*MetadataListingCacheAdapter)(nil)

type MetadataListingCacheAdapter struct {
	cache cache.Cache[file.Listing]
}

func NewMetadataListingCacheAdapter(cache cache.Cache[file.Listing]) *MetadataListingCacheAdapter {
	return &MetadataListingCacheAdapter{cache: cache}
}

func (a *MetadataListingCacheAdapter) Get(ctx context.Context, path string) (*file.Listing, bool) {
	listing, err := a.cache.Get(ctx, path)
	if err != nil {
		return nil, false
	}
	return &listing, true
}

func (a *MetadataListingCacheAdapter) Put(ctx context.Context, path string, listing *file.Listing) error {
	return a.cache.Set(ctx, path, *listing)
}

func (a *MetadataListingCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return a.cache.Delete(ctx, path)
}
