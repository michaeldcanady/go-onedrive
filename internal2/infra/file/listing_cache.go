package file

import "context"

// ListingCache defines the interface for caching lists of drive items.
// It maps a path to a Listing, which contains the IDs of the children.
type ListingCache interface {
	// Get retrieves a cached listing for the given path.
	Get(ctx context.Context, path string) (*Listing, bool)
	// Put stores a listing in the cache for the given path.
	Put(ctx context.Context, path string, l *Listing) error
	// Invalidate removes a cached listing for the given path.
	Invalidate(ctx context.Context, path string) error
}
