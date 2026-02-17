package file

import "context"

type Listing struct {
	ETag     string
	ChildIDs []string
}

type ListingCache interface {
	Get(ctx context.Context, path string) (*Listing, bool)
	Put(ctx context.Context, path string, listing *Listing) error
	Invalidate(ctx context.Context, path string) error
}
