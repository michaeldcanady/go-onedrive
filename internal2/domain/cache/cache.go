package cache

import (
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

// Alias types from pkg/cache to maintain domain layer abstraction.
type (
	Cache[T any]                  = pkgcache.Cache[T]
	Entry[K comparable, V any]    = pkgcache.Entry[K, V]
	SerializerFunc                = pkgcache.SerializerFunc
	DeserializerFunc              = pkgcache.DeserializerFunc
	Serializer[T any]             = pkgcache.Serializer[T]
	Deserializer[T any]           = pkgcache.Deserializer[T]
	SerializerDeserializer[T any] = pkgcache.SerializerDeserializer[T]
	KeyValueStore                 = pkgcache.KeyValueStore
	Store                         = pkgcache.Store
)

// NewEntry wraps the pkg/cache NewEntry function.
func NewEntry[K comparable, V any](key K, value V) *Entry[K, V] {
	return pkgcache.NewEntry(key, value)
}

// NewStore wraps the pkg/cache NewStore function.
func NewStore(store KeyValueStore) *Store {
	return pkgcache.NewStore(store)
}
