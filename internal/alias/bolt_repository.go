package alias

import (
	"context"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type BoltRepository struct {
	db *bolt.DB
}

func NewBoltRepository(db *bolt.DB) *BoltRepository {
	return &BoltRepository{db: db}
}

func (r *BoltRepository) getBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte("drive_aliases"))
	if b == nil {
		return nil, fmt.Errorf("drive_aliases bucket not found")
	}
	return b, nil
}

func (r *BoltRepository) Get(ctx context.Context, name string) (string, error) {
	var driveID string
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getBucket(tx)
		if err != nil {
			return err
		}
		data := b.Get([]byte(name))
		if data == nil {
			return fmt.Errorf("alias not found")
		}
		driveID = string(data)
		return nil
	})
	return driveID, err
}

func (r *BoltRepository) Set(ctx context.Context, name, driveID string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := r.getBucket(tx)
		if err != nil {
			return err
		}
		return b.Put([]byte(name), []byte(driveID))
	})
}

func (r *BoltRepository) Delete(ctx context.Context, name string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := r.getBucket(tx)
		if err != nil {
			return err
		}
		return b.Delete([]byte(name))
	})
}

func (r *BoltRepository) List(ctx context.Context) (map[string]string, error) {
	aliases := make(map[string]string)
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getBucket(tx)
		if err != nil {
			return err
		}
		return b.ForEach(func(k, v []byte) error {
			aliases[string(k)] = string(v)
			return nil
		})
	})
	return aliases, err
}
