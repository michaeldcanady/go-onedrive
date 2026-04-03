package state

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/profile"

	bolt "go.etcd.io/bbolt"
)

var (
	// globalBucketName is used for state that should persist across application restarts and be shared across sessions.
	globalBucketName = []byte("global")
	// sessionBucketName is used for temporary state that should not persist across application restarts.
	sessionBucketName = []byte("session")
	// driveAliasesBucketName is used to store user-defined drive aliases.
	driveAliasesBucketName = []byte("drive_aliases")
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
		if err := bs.Set(KeyProfile, profile.DefaultProfileName, ScopeGlobal); err != nil {
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
		if _, err := tx.CreateBucketIfNotExists(driveAliasesBucketName); err != nil {
			return fmt.Errorf("failed to create drive_aliases bucket: %w", err)
		}
		return nil
	})
}

// Get retrieves a state value by its key, checking session scope first, then global.
func (bs *BoltService) Get(key Key) (string, error) {
	var value string
	err := bs.db.View(func(tx *bolt.Tx) error {
		keyStr := key.String()

		for _, scopeKey := range [][]byte{sessionBucketName, globalBucketName} {
			b := tx.Bucket(scopeKey)
			if b == nil {
				return ErrBucketNotFound
			}

			v := b.Get([]byte(keyStr))
			if v == nil {
				continue // Key not found in this scope, try the next one.
			}
			if value = string(v); value != "" {
				return nil // Return immediately if a non-empty value is found
			}
		}
		return ErrKeyNotFound
	})

	if err != nil {
		return "", err
	}
	return value, nil
}

// Set assigns a value to a key within the specified scope.
func (bs *BoltService) Set(key Key, value string, scope Scope) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		keyStr := key.String()
		bucketName := globalBucketName
		if scope == ScopeSession {
			bucketName = sessionBucketName
		}

		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("failed to get or create bucket %s: %w", string(bucketName), err)
		}

		if err := b.Put([]byte(keyStr), []byte(value)); err != nil {
			return fmt.Errorf("failed to put state key %s: %w", keyStr, err)
		}
		return nil
	})
}

// Clear removes a state value for the given key from all scopes.
func (bs *BoltService) Clear(key Key) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		keyStr := key.String()

		for _, bucketName := range [][]byte{sessionBucketName, globalBucketName} {
			b := tx.Bucket(bucketName)
			if b == nil {
				return fmt.Errorf("bucket %s not found", string(bucketName))
			}
			if err := b.Delete([]byte(keyStr)); err != nil {
				return fmt.Errorf("failed to delete state key %s from bucket %s: %w", keyStr, string(bucketName), err)
			}
		}
		return nil
	})
}

// TODO: drive alias management should likely be moved to its own service

// GetDriveAlias retrieves the drive ID associated with a user-defined alias.
func (bs *BoltService) GetDriveAlias(alias string) (string, error) {
	var value string
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(driveAliasesBucketName)
		if b == nil {
			return fmt.Errorf("drive_aliases bucket not found")
		}
		v := b.Get([]byte(alias))
		if v == nil {
			return fmt.Errorf("drive alias %s not found", alias)
		}
		value = string(v)
		return nil
	})
	return value, err
}

// SetDriveAlias assigns a drive ID to a specified alias.
func (bs *BoltService) SetDriveAlias(alias, driveID string) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(driveAliasesBucketName)
		if err != nil {
			return fmt.Errorf("failed to get or create drive_aliases bucket: %w", err)
		}
		if err := b.Put([]byte(alias), []byte(driveID)); err != nil {
			return fmt.Errorf("failed to put drive alias %s: %w", alias, err)
		}
		return nil
	})
}

// RemoveDriveAlias deletes a drive alias.
func (bs *BoltService) RemoveDriveAlias(alias string) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(driveAliasesBucketName)
		if b == nil {
			return fmt.Errorf("drive_aliases bucket not found")
		}
		if err := b.Delete([]byte(alias)); err != nil {
			return fmt.Errorf("failed to delete drive alias %s: %w", alias, err)
		}
		return nil
	})
}

// ListDriveAliases returns all registered drive aliases.
func (bs *BoltService) ListDriveAliases() (map[string]string, error) {
	aliases := make(map[string]string)
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(driveAliasesBucketName)
		if b == nil {
			return fmt.Errorf("drive_aliases bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			aliases[string(k)] = string(v)
			return nil
		})
	})
	return aliases, err
}

func (bs *BoltService) GetDriveAliasByDriveID(driveID string) (string, error) {
	var alias string
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(driveAliasesBucketName)
		if b == nil {
			return fmt.Errorf("drive_aliases bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			if string(v) == driveID {
				alias = string(k)
				return nil // Stop iteration once we find the matching drive ID
			}
			return nil
		})
	})
	if err != nil {
		return "", err
	}
	if alias == "" {
		return "", nil
	}
	return alias, nil
}
