package cache

import (
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
)

// BoltCacheFactory returns a factory function that creates a new BoltDB-backed KeyValueStore.
func BoltCacheFactory(path, bucket string) func() abstractions.KeyValueStore {
	return func() abstractions.KeyValueStore {
		store, err := bolt.NewStore(path, bucket)
		if err != nil {
			return nil
		}
		return store
	}
}

// SiblingBoltFactory returns a factory function that creates a new BoltDB-backed KeyValueStore
// using an existing BoltDB instance but a different bucket.
func SiblingBoltFactory(store *bolt.Store, bucket string) func() abstractions.KeyValueStore {
	return func() abstractions.KeyValueStore {
		siblingStore, err := bolt.NewSiblingStore(store, bucket)
		if err != nil {
			return nil
		}
		return siblingStore
	}
}
