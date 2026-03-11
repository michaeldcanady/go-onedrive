package infra

import (
	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/michaeldcanady/go-onedrive/pkg/cache/sqlite"
)

type Cache[K comparable, V any] = sqlite.Cache[K, V]

func New[K comparable, V any](
	path string,
	keySer cachedomain.SerializerDeserializer[K],
	valueSer cachedomain.SerializerDeserializer[V],
) (*Cache[K, V], error) {
	return sqlite.New(path, keySer, valueSer)
}
