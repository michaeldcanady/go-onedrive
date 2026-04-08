package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/shared"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	bolt "go.etcd.io/bbolt"
)

// BoltService is a persistent implementation of the profile.Service using BoltDB.
type BoltService struct {
	db    *bolt.DB
	env   environment.Service
	state state.Service
	log   logger.Logger
}

// NewBoltService initializes a new instance of the BoltService.
func NewBoltService(env environment.Service, state state.Service, l logger.Logger) (*BoltService, error) {
	configDir, err := env.ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	dbFilePath := filepath.Join(configDir, "profiles.db")
	l.Debug("opening profiles database", logger.String("path", dbFilePath))
	db, err := bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	bs := &BoltService{
		db:    db,
		env:   env,
		state: state,
		log:   l,
	}

	// Ensure profiles bucket is created
	if err := bs.ensureBucket(); err != nil {
		bs.db.Close() // Close DB if initialization fails
		return nil, err
	}

	// Ensure default profile exists
	_, _ = bs.Create(context.Background(), shared.DefaultProfileName)

	if err := bs.migrateConfigPaths(); err != nil {
		bs.db.Close()
		return nil, fmt.Errorf("failed to migrate profile config paths: %w", err)
	}

	return bs, nil
}

// ensureBuckets creates the top-level buckets if they don't exist.
func (bs *BoltService) ensureBucket() error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("profiles")); err != nil {
			return fmt.Errorf("failed to create profiles bucket: %w", err)
		}
		return nil
	})
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
				bs.log.Warn("skipping invalid profile entry during migration", logger.String("key", string(k)))
				return nil // Skip invalid entries
			}

			changed := false
			if p.ConfigPath == "" {
				p.ConfigPath = filepath.Join(configDir, fmt.Sprintf("%s.yaml", p.Name))
				changed = true
			}

			if changed {
				bs.log.Info("migrating profile config path", logger.String("profile", p.Name), logger.String("new_path", p.ConfigPath))
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
	bs.log.WithContext(ctx).Debug("resolving profile config path", logger.String("profile", profileName))
	p, err := bs.Get(ctx, profileName)
	if err != nil {
		return "", err
	}
	return p.ConfigPath, nil
}

// Close closes the BoltDB database connection.
func (bs *BoltService) Close() error {
	bs.log.Debug("closing profiles database")
	return bs.db.Close()
}

func (bs *BoltService) getBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte("profiles"))
	if b == nil {
		bs.log.Error("profiles bucket not found", logger.String("bucket", "profiles"))
		return nil, ErrProfilesBucketNotFound
	}
	return b, nil
}

// Get returns the profile with the specified name.
func (bs *BoltService) Get(ctx context.Context, name string) (Profile, error) {
	bs.log.WithContext(ctx).Debug("retrieving profile", logger.String("profile", name))
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
	bs.log.WithContext(ctx).Debug("listing profiles")
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
	log := bs.log.WithContext(ctx).With(logger.String("profile", name))
	log.Info("creating new profile")

	configDir, err := bs.env.ConfigDir()
	if err != nil {
		log.Error("failed to get config directory", logger.Error(err))
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
	if err != nil {
		if !errors.Is(err, ErrProfileAlreadyExists) {
			log.Error("failed to create profile", logger.Error(err))
		} else {
			log.Warn("profile already exists")
		}
		return Profile{}, err
	}
	return p, nil
}

// Update saves the specified profile.
func (bs *BoltService) Update(ctx context.Context, p Profile) error {
	log := bs.log.WithContext(ctx).With(logger.String("profile", p.Name))
	log.Info("updating profile")

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
	log := bs.log.WithContext(ctx).With(logger.String("profile", name))
	log.Info("deleting profile")

	if name == shared.DefaultProfileName {
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

// GetActive retrieves the currently active profile.
func (bs *BoltService) GetActive(ctx context.Context) (Profile, error) {
	name, err := bs.state.Get(state.KeyProfile)
	if err != nil {
		bs.log.WithContext(ctx).Error("failed to get active profile name from state", logger.Error(err))
		return Profile{}, fmt.Errorf("failed to get active profile name: %w", err)
	}

	bs.log.WithContext(ctx).Debug("retrieved active profile", logger.String("profile", name))
	return bs.Get(ctx, name)
}

// SetActive marks a specific profile as the active one with the given scope.
func (bs *BoltService) SetActive(ctx context.Context, name string, scope state.Scope) error {
	log := bs.log.WithContext(ctx).With(logger.String("profile", name), logger.String("scope", scope.String()))
	log.Info("setting active profile")

	exists, err := bs.Exists(ctx, name)
	if err != nil {
		log.Error("failed to check profile existence during SetActive", logger.Error(err))
		return err
	}
	if !exists {
		log.Warn("profile not found during SetActive")
		return ErrProfileNotFound
	}

	return bs.state.Set(state.KeyProfile, name, scope)
}
