package state

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
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
	// driveAliasesBucketName is used to store user-defined drive aliases.
	driveAliasesBucketName = []byte("drive_aliases")
)

// TODO: consolidate DefaultProfileName somewhere
const (
	// stateDBFileName is the name of the BoltDB file used for storing state data.
	stateDBFileName = "state.db"
)

// BoltService is a persistent implementation of the state.Service using BoltDB for global state
// and an in-memory map for session state.
type BoltService struct {
	db           *bolt.DB
	sessionState map[Key]string
	mu           sync.RWMutex
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
		db:           db,
		sessionState: make(map[Key]string),
	}

	// Ensure top-level buckets are created
	if err := bs.ensureBuckets(); err != nil {
		bs.db.Close() // Close DB if initialization fails
		return nil, err
	}

	// Set default profile if not set in persistent storage
	if _, err := bs.getGlobal(KeyProfile); err != nil {
		if err != ErrKeyNotFound {
			return nil, fmt.Errorf("failed to check for default profile: %w", err)
		}
		if err := bs.setGlobal(KeyProfile, shared.DefaultProfileName); err != nil {
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
		return nil
	})
}

// Get retrieves a state value by its key, checking session scope first, then global.
func (bs *BoltService) Get(key Key) (string, error) {
	// Check session state first
	bs.mu.RLock()
	val, ok := bs.sessionState[key]
	bs.mu.RUnlock()

	if ok && val != "" {
		return val, nil
	}

	// Check global state
	return bs.getGlobal(key)
}

func (bs *BoltService) getGlobal(key Key) (string, error) {
	var value string
	err := bs.db.View(func(tx *bolt.Tx) error {
		keyStr := key.String()
		b := tx.Bucket(globalBucketName)
		if b == nil {
			return ErrBucketNotFound
		}

		v := b.Get([]byte(keyStr))
		if v == nil {
			return ErrKeyNotFound
		}
		value = string(v)
		return nil
	})

	if err != nil {
		return "", err
	}
	return value, nil
}

// Set assigns a value to a key within the specified scope.
func (bs *BoltService) Set(key Key, value string, scope Scope) error {
	if scope == ScopeSession {
		bs.mu.Lock()
		bs.sessionState[key] = value
		bs.mu.Unlock()
		return nil
	}

	return bs.setGlobal(key, value)
}

func (bs *BoltService) setGlobal(key Key, value string) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		keyStr := key.String()
		b, err := tx.CreateBucketIfNotExists(globalBucketName)
		if err != nil {
			return fmt.Errorf("failed to get or create bucket %s: %w", string(globalBucketName), err)
		}

		if err := b.Put([]byte(keyStr), []byte(value)); err != nil {
			return fmt.Errorf("failed to put state key %s: %w", keyStr, err)
		}
		return nil
	})
}

// Clear removes a state value for the given key from all scopes.
func (bs *BoltService) Clear(key Key) error {
	bs.mu.Lock()
	delete(bs.sessionState, key)
	bs.mu.Unlock()

	return bs.db.Update(func(tx *bolt.Tx) error {
		keyStr := key.String()
		b := tx.Bucket(globalBucketName)
		if b == nil {
			return nil // No bucket, nothing to clear
		}
		if err := b.Delete([]byte(keyStr)); err != nil {
			return fmt.Errorf("failed to delete state key %s from global bucket: %w", keyStr, err)
		}
		return nil
	})
}

// GetProfile retrieves the currently active profile name.
func (bs *BoltService) GetProfile(_ context.Context) (string, error) {
	return bs.Get(KeyProfile)
}

// SetProfile updates the active profile name.
func (bs *BoltService) SetProfile(_ context.Context, name string, scope Scope) error {
	return bs.Set(KeyProfile, name, scope)
}

// GetDrive retrieves the active drive ID.
func (bs *BoltService) GetDrive(_ context.Context) (string, error) {
	return bs.Get(KeyDrive)
}

// SetDrive updates the active drive ID.
func (bs *BoltService) SetDrive(_ context.Context, driveID string, scope Scope) error {
	return bs.Set(KeyDrive, driveID, scope)
}

// GetAccessToken retrieves the cached access token.
func (bs *BoltService) GetAccessToken(_ context.Context) (string, error) {
	return bs.Get(KeyAccessToken)
}

// SetAccessToken updates the cached access token.
func (bs *BoltService) SetAccessToken(_ context.Context, token string, scope Scope) error {
	return bs.Set(KeyAccessToken, token, scope)
}

// ClearAccessToken removes the cached access token.
func (bs *BoltService) ClearAccessToken(_ context.Context) error {
	return bs.Clear(KeyAccessToken)
}

// GetConfigOverride retrieves the configuration path override.
func (bs *BoltService) GetConfigOverride(_ context.Context) (string, error) {
	return bs.Get(KeyConfigOverride)
}

// SetConfigOverride updates the configuration path override.
func (bs *BoltService) SetConfigOverride(_ context.Context, path string, scope Scope) error {
	return bs.Set(KeyConfigOverride, path, scope)
}
