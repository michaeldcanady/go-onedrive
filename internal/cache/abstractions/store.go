package abstractions

import "context"

// Store is a generic blob store keyed by a logical name (e.g., "profile").
type Store interface {
	// LoadBytes returns the raw bytes for the given key, or (nil, nil) if not found.
	LoadBytes(ctx context.Context, key string) ([]byte, error)

	// SaveBytes persists the raw bytes for the given key.
	SaveBytes(ctx context.Context, key string, data []byte) error

	// Delete removes the data for the given key.
	Delete(ctx context.Context, key string) error
}
