package disk

import (
	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/michaeldcanady/go-onedrive/pkg/cache/disk"
)

type Cache[K comparable, V any] = disk.Cache[K, V]

func New[K comparable, V any](
	path string,
	ks cache.SerializerDeserializer[K],
	vs cache.SerializerDeserializer[V],
) (*Cache[K, V], error) {
	return disk.New(path, ks, vs)
}
