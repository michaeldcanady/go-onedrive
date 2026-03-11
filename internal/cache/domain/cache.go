package domain

import (
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

// Alias types from pkg/cache to maintain domain layer abstraction.
type (
	Cache[T any]                  = pkgcache.Cache[T]
	SerializerFunc                = pkgcache.SerializerFunc
	DeserializerFunc              = pkgcache.DeserializerFunc
	Serializer[T any]             = pkgcache.Serializer[T]
	Deserializer[T any]           = pkgcache.Deserializer[T]
	SerializerDeserializer[T any] = pkgcache.SerializerDeserializer[T]
	KeyValueStore                 = pkgcache.KeyValueStore
	Store                         = pkgcache.Store
)

// NewStore wraps the pkg/cache NewStore function.
func NewStore(store KeyValueStore) *Store {
	return pkgcache.NewStore(store)
}
