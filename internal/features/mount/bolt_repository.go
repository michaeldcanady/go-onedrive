package mount

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	mountsBucket = []byte("mounts")
)

type boltRepository struct {
	db *bbolt.DB
}

// NewBoltRepository creates a new bbolt-based mount repository.
func NewBoltRepository(db *bbolt.DB) (Repository, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(mountsBucket)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mount bucket: %w", err)
	}

	return &boltRepository{db: db}, nil
}

func (r *boltRepository) Save(m *Mount) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mountsBucket)
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return b.Put([]byte(m.Path), data)
	})
}

func (r *boltRepository) List() ([]*Mount, error) {
	var mounts []*Mount
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mountsBucket)
		return b.ForEach(func(k, v []byte) error {
			var m Mount
			if err := json.Unmarshal(v, &m); err != nil {
				return err
			}
			mounts = append(mounts, &m)
			return nil
		})
	})
	return mounts, err
}

func (r *boltRepository) Delete(path string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mountsBucket)
		return b.Delete([]byte(path))
	})
}

func (r *boltRepository) Get(path string) (*Mount, error) {
	var m Mount
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mountsBucket)
		v := b.Get([]byte(path))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &m)
	})
	if m.Path == "" {
		return nil, nil
	}
	return &m, err
}
