package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// BoltRepository implements both ProfileRepository and SettingsRepository using BoltDB.
type BoltRepository struct {
	db *bolt.DB
}

// NewBoltRepository creates a new instance of BoltRepository and initializes the DB schema.
func NewBoltRepository(dbPath string) (*BoltRepository, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	// Ensure buckets are created
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("profiles")); err != nil {
			return fmt.Errorf("failed to create profiles bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("settings")); err != nil {
			return fmt.Errorf("failed to create settings bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &BoltRepository{db: db}, nil
}

// Close closes the underlying BoltDB database.
func (r *BoltRepository) Close() error {
	return r.db.Close()
}

func (r *BoltRepository) getProfileBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte("profiles"))
	if b == nil {
		return nil, ErrProfilesBucketNotFound
	}
	return b, nil
}

func (r *BoltRepository) getSettingsBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte("settings"))
	if b == nil {
		return nil, fmt.Errorf("settings bucket not found")
	}
	return b, nil
}

// --- ProfileRepository Implementation ---

func (r *BoltRepository) Get(ctx context.Context, name string) (Profile, error) {
	var p Profile
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getProfileBucket(tx)
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

func (r *BoltRepository) Create(ctx context.Context, p Profile) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := r.getProfileBucket(tx)
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

func (r *BoltRepository) Update(ctx context.Context, p Profile) error {
	return r.Create(ctx, p)
}

func (r *BoltRepository) Delete(ctx context.Context, name string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := r.getProfileBucket(tx)
		if err != nil {
			return err
		}
		return b.Delete([]byte(name))
	})
}

func (r *BoltRepository) List(ctx context.Context) ([]Profile, error) {
	var profiles []Profile
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getProfileBucket(tx)
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

func (r *BoltRepository) Exists(ctx context.Context, name string) (bool, error) {
	exists := false
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getProfileBucket(tx)
		if err != nil {
			return err
		}
		if b.Get([]byte(name)) != nil {
			exists = true
		}
		return nil
	})
	return exists, err
}

// --- SettingsRepository Implementation ---

func (r *BoltRepository) GetSetting(ctx context.Context, key string) (string, error) {
	var value string
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getSettingsBucket(tx)
		if err != nil {
			return err
		}
		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("setting not found: %s", key)
		}
		value = string(data)
		return nil
	})
	return value, err
}

func (r *BoltRepository) SetSetting(ctx context.Context, key, value string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("settings"))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
}
