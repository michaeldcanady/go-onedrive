package app

import (
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/michaeldcanady/go-onedrive/pkg/cache/bolt"
)

// BoltCacheFactory returns a factory function that creates a new BoltDB-backed KeyValueStore.
func BoltCacheFactory(path, bucket string) func() pkgcache.KeyValueStore {
	return func() pkgcache.KeyValueStore {
		store, err := bolt.NewStore(path, bucket)
		if err != nil {
			return nil
		}
		return store
	}
}

// SiblingBoltFactory returns a factory function that creates a new BoltDB-backed KeyValueStore
// using an existing BoltDB instance but a different bucket.
func SiblingBoltFactory(store *bolt.Store, bucket string) func() pkgcache.KeyValueStore {
	return func() pkgcache.KeyValueStore {
		siblingStore, err := bolt.NewSiblingStore(store, bucket)
		if err != nil {
			return nil
		}
		return siblingStore
	}
}
