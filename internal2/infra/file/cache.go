package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

// cache is an internal interface defining the behavior for interacting with
// a generic cache. It abstracts the underlying storage mechanism and
// serialization/deserialization logic.
type cache interface {
	// Delete removes an entry from the cache identified by the keySerializer.
	Delete(ctx context.Context, keySerializer abstractions.Serializer2) error
	// Get retrieves an entry from the cache identified by the keySerializer and
	// populates the valueDeserializer.
	Get(ctx context.Context, keySerializer abstractions.Serializer2, valueDeserializer abstractions.Deserializer2) error
	// Set stores an entry in the cache identified by the keySerializer and
	// serialized by the valueSerializer.
	Set(ctx context.Context, keySerializer abstractions.Serializer2, valueSerializer abstractions.Serializer2) error
}
