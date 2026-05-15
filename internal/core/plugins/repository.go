package plugins

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	bucketName = []byte("plugins")
)

// Repository manages the persistence of plugin metadata.
type Repository interface {
	Get(path string) (*Metadata, error)
	Set(path string, meta *Metadata) error
	Delete(path string) error
	List() ([]*Metadata, error)
}

type boltRepository struct {
	db *bbolt.DB
}

// NewBoltRepository returns a new [Repository] backed by a bbolt database.
func NewBoltRepository(db *bbolt.DB) (Repository, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create plugins bucket: %w", err)
	}

	return &boltRepository{db: db}, nil
}

func (r *boltRepository) Get(path string) (*Metadata, error) {
	var meta *Metadata
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		v := b.Get([]byte(path))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &meta)
	})
	return meta, err
}

func (r *boltRepository) Set(path string, meta *Metadata) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		v, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		return b.Put([]byte(path), v)
	})
}

func (r *boltRepository) Delete(path string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Delete([]byte(path))
	})
}

func (r *boltRepository) List() ([]*Metadata, error) {
	var results []*Metadata
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.ForEach(func(k, v []byte) error {
			var meta Metadata
			if err := json.Unmarshal(v, &meta); err != nil {
				return err
			}
			results = append(results, &meta)
			return nil
		})
	})
	return results, err
}
