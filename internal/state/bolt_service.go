package state

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/shared"

	bolt "go.etcd.io/bbolt"
)

var (
	// globalBucketName is used for state that should persist across application restarts and be shared across sessions.
	globalBucketName = []byte("global")
	// sessionBucketName is used for temporary state that should not persist across application restarts.
	sessionBucketName = []byte("session")
)

const (
	// stateDBFileName is the name of the BoltDB file used for storing state data.
	stateDBFileName = "state.db"
)

// BoltService is a persistent implementation of the state.Service using BoltDB.
type BoltService struct {
	db *bolt.DB
}

// NewBoltService initializes a new instance of the BoltService.
func NewBoltService(env environment.Service) (*BoltService, error) {
	dbPath, err := env.StateDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get state directory: %w", err)
	}

	dbFilePath := filepath.Join(dbPath, stateDBFileName)
	db, err := bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	bs := &BoltService{
		db: db,
	}

	// Ensure top-level buckets are created
	if err := bs.ensureBuckets(); err != nil {
		bs.db.Close() // Close DB if initialization fails
		return nil, err
	}

	// Set default profile if not set
	if _, err := bs.Get(KeyProfile); err != nil {
		if err != ErrKeyNotFound {
			return nil, fmt.Errorf("failed to check for default profile: %w", err)
		}
		if err := bs.Set(KeyProfile, shared.DefaultProfileName, ScopeGlobal); err != nil {
			return nil, fmt.Errorf("failed to set default profile: %w", err)
		}
	}

	return bs, nil
}

// Close closes the BoltDB database connection.
func (bs *BoltService) Close() error {
	return bs.db.Close()
}

// ensureBuckets creates the top-level buckets if they don't exist.
func (bs *BoltService) ensureBuckets() error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(globalBucketName); err != nil {
			return fmt.Errorf("failed to create global bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(sessionBucketName); err != nil {
			return fmt.Errorf("failed to create session bucket: %w", err)
		}
		return nil
	})
}

// Get retrieves a state value by its key, checking session scope first, then global.
func (bs *BoltService) Get(key Key) (string, error) {
	return bs.GetScoped("", key.String())
}

// Set assigns a value to a key within the specified scope.
func (bs *BoltService) Set(key Key, value string, scope Scope) error {
	return bs.SetScoped("", key.String(), value, scope)
}

// Clear removes a state value for the given key from all scopes.
func (bs *BoltService) Clear(key Key) error {
	return bs.ClearScoped("", key.String())
}

// GetScoped retrieves a value from a named sub-bucket.
func (bs *BoltService) GetScoped(bucket, key string) (string, error) {
	var value string
	err := bs.db.View(func(tx *bolt.Tx) error {
		for _, scopeKey := range [][]byte{sessionBucketName, globalBucketName} {
			root := tx.Bucket(scopeKey)
			if root == nil {
				continue
			}

			b := root
			if bucket != "" {
				b = root.Bucket([]byte(bucket))
				if b == nil {
					continue
				}
			}

			v := b.Get([]byte(key))
			if v != nil {
				value = string(v)
				return nil
			}
		}
		return ErrKeyNotFound
	})

	return value, err
}

// SetScoped assigns a value to a key within a named sub-bucket.
func (bs *BoltService) SetScoped(bucket, key, value string, scope Scope) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		rootName := globalBucketName
		if scope == ScopeSession {
			rootName = sessionBucketName
		}

		root, err := tx.CreateBucketIfNotExists(rootName)
		if err != nil {
			return err
		}

		b := root
		if bucket != "" {
			b, err = root.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}

		return b.Put([]byte(key), []byte(value))
	})
}

// ClearScoped removes a value from a named sub-bucket.
func (bs *BoltService) ClearScoped(bucket, key string) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		for _, rootName := range [][]byte{sessionBucketName, globalBucketName} {
			root := tx.Bucket(rootName)
			if root == nil {
				continue
			}

			b := root
			if bucket != "" {
				b = root.Bucket([]byte(bucket))
				if b == nil {
					continue
				}
			}

			if err := b.Delete([]byte(key)); err != nil {
				return err
			}
		}
		return nil
	})
}

// ListScoped returns all keys within a named sub-bucket across all scopes.
func (bs *BoltService) ListScoped(bucket string) ([]string, error) {
	keySet := make(map[string]struct{})
	err := bs.db.View(func(tx *bolt.Tx) error {
		for _, rootName := range [][]byte{sessionBucketName, globalBucketName} {
			root := tx.Bucket(rootName)
			if root == nil {
				continue
			}

			b := root
			if bucket != "" {
				b = root.Bucket([]byte(bucket))
				if b == nil {
					continue
				}
			}

			err := b.ForEach(func(k, v []byte) error {
				// If v is nil, it's a sub-bucket, which we ignore for now as we only want keys in the current bucket level.
				if v != nil {
					keySet[string(k)] = struct{}{}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}

	return keys, nil
}
