package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

type cache interface {
	Delete(ctx context.Context, keySerializer abstractions.Serializer2) error
	Get(ctx context.Context, keySerializer abstractions.Serializer2, valueDeserializer abstractions.Deserializer2) error
	Set(ctx context.Context, keySerializer abstractions.Serializer2, valueSerializer abstractions.Serializer2) error
}
