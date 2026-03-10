package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

// ListingCache defines the interface for a cache that stores directory listings.
type ListingCache interface {
	Get(ctx context.Context, path string) (*file.Listing, bool)
	Put(ctx context.Context, path string, listing *file.Listing) error
	Invalidate(ctx context.Context, path string) error
}
