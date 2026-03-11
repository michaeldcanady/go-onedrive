package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// ListingCache defines the interface for a cache that stores directory listings.
type ListingCache interface {
	Get(ctx context.Context, path string) (*domain.Listing, bool)
	Put(ctx context.Context, path string, listing *domain.Listing) error
	Invalidate(ctx context.Context, path string) error
}
