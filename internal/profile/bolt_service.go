package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	bolt "go.etcd.io/bbolt"
)

// BoltService is a persistent implementation of the profile.Service using BoltDB.
type BoltService struct {
	db  *bolt.DB
	env environment.Service
}

// NewBoltService initializes a new instance of the BoltService.
func NewBoltService(env environment.Service) (*BoltService, error) {
	configDir, err := env.ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	dbFilePath := filepath.Join(configDir, "profiles.db")
	db, err := bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	bs := &BoltService{
		db:  db,
		env: env,
	}

	// Ensure profiles bucket is created
	if err := bs.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("profiles"))
		return err
	}); err != nil {
		bs.db.Close()
		return nil, fmt.Errorf("failed to create profiles bucket: %w", err)
	}

	// Ensure default profile exists
	_, _ = bs.Create(context.Background(), DefaultProfileName)

	if err := bs.migrateConfigPaths(); err != nil {
		bs.db.Close()
		return nil, fmt.Errorf("failed to migrate profile config paths: %w", err)
	}

	return bs, nil
}

func (bs *BoltService) migrateConfigPaths() error {
	configDir, err := bs.env.ConfigDir()
	if err != nil {
		return err
	}

	return bs.db.Update(func(tx *bolt.Tx) error {
		b, err := bs.getBucket(tx)
		if err != nil {
			return err
		}

		return b.ForEach(func(k, v []byte) error {
			var p Profile
			if err := json.Unmarshal(v, &p); err != nil {
				return nil // Skip invalid entries
			}

			changed := false
			if p.ConfigPath == "" {
				p.ConfigPath = filepath.Join(configDir, fmt.Sprintf("%s.yaml", p.Name))
				changed = true
			}

			if changed {
				data, err := json.Marshal(p)
				if err != nil {
					return err
				}
				return b.Put(k, data)
			}
			return nil
		})
	})
}

// ResolvePath returns the configuration file path for the specified profile name.
func (bs *BoltService) ResolvePath(ctx context.Context, profileName string) (string, error) {
	p, err := bs.Get(ctx, profileName)
	if err != nil {
		return "", err
	}
	return p.ConfigPath, nil
}

// Close closes the BoltDB database connection.
func (bs *BoltService) Close() error {
	return bs.db.Close()
}

func (bs *BoltService) getBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte("profiles"))
	if b == nil {
		return nil, ErrProfilesBucketNotFound
	}
	return b, nil
}

// Get returns the profile with the specified name.
func (bs *BoltService) Get(ctx context.Context, name string) (Profile, error) {
	var p Profile
	err := bs.db.View(func(tx *bolt.Tx) error {
		b, err := bs.getBucket(tx)
		if err != nil {
			return err
		}
		data := b.Get([]byte(name))
		if data == nil {
			return ErrProfileNotFound
		}
		return json.Unmarshal(data, &p)
	})
	return p, err
}

// List returns a list of all profiles.
func (bs *BoltService) List(ctx context.Context) ([]Profile, error) {
	var profiles []Profile
	err := bs.db.View(func(tx *bolt.Tx) error {
		b, err := bs.getBucket(tx)
		if err != nil {
			return err
		}
		return b.ForEach(func(k, v []byte) error {
			var p Profile
			if err := json.Unmarshal(v, &p); err == nil {
				profiles = append(profiles, p)
			}
			return nil
		})
	})
	return profiles, err
}

// Create generates a new profile with the specified name.
func (bs *BoltService) Create(ctx context.Context, name string) (Profile, error) {
	configDir, err := bs.env.ConfigDir()
	if err != nil {
		return Profile{}, fmt.Errorf("failed to get config directory: %w", err)
	}

	p := Profile{
		Name:       name,
		ConfigPath: filepath.Join(configDir, fmt.Sprintf("%s.yaml", name)),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b, err := bs.getBucket(tx)
		if err != nil {
			return err
		}
		if b.Get([]byte(name)) != nil {
			return ErrProfileAlreadyExists
		}
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(name), data)
	})
	return p, err
}

// Update saves the specified profile.
func (bs *BoltService) Update(ctx context.Context, p Profile) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		b, err := bs.getBucket(tx)
		if err != nil {
			return err
		}
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(p.Name), data)
	})
}

// Delete removes the specified profile name.
func (bs *BoltService) Delete(ctx context.Context, name string) error {
	if name == DefaultProfileName {
		return fmt.Errorf("cannot delete the default profile")
	}
	return bs.db.Update(func(tx *bolt.Tx) error {
		b, err := bs.getBucket(tx)
		if err != nil {
			return err
		}
		return b.Delete([]byte(name))
	})
}

// Exists checks if a profile with the specified name exists.
func (bs *BoltService) Exists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := bs.db.View(func(tx *bolt.Tx) error {
		b, err := bs.getBucket(tx)
		if err != nil {
			return err
		}
		exists = b.Get([]byte(name)) != nil
		return nil
	})
	return exists, err
}
