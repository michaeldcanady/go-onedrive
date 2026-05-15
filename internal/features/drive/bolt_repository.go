package drive

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	drivesBucket = []byte("drives")
)

type boltRepository struct {
	db *bbolt.DB
}

// NewBoltRepository creates a new bbolt-based drive repository.
func NewBoltRepository(db *bbolt.DB) (Repository, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(drivesBucket)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize drives bucket: %w", err)
	}

	return &boltRepository{db: db}, nil
}

func (r *boltRepository) Save(d *Drive) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(drivesBucket)
		data, err := json.Marshal(d)
		if err != nil {
			return err
		}
		return b.Put([]byte(d.ID), data)
	})
}

func (r *boltRepository) ListByIdentity(identityID string) ([]*Drive, error) {
	var drives []*Drive
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(drivesBucket)
		return b.ForEach(func(k, v []byte) error {
			var d Drive
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			if identityID == "" || d.IdentityID == identityID {
				drives = append(drives, &d)
			}
			return nil
		})
	})
	return drives, err
}

func (r *boltRepository) ByID(driveID string) (*Drive, error) {
	var d Drive
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(drivesBucket)
		v := b.Get([]byte(driveID))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &d)
	})
	if d.ID == "" {
		return nil, nil
	}
	return &d, err
}

func (r *boltRepository) Delete(driveID string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(drivesBucket)
		return b.Delete([]byte(driveID))
	})
}
