package profile

import (
	"context"
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

type BoltRepository struct {
	db *bolt.DB
}

func NewBoltRepository(db *bolt.DB) *BoltRepository {
	return &BoltRepository{db: db}
}

func (r *BoltRepository) getBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte("profiles"))
	if b == nil {
		return nil, ErrProfilesBucketNotFound
	}
	return b, nil
}

func (r *BoltRepository) Get(ctx context.Context, name string) (Profile, error) {
	var p Profile
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getBucket(tx)
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
		b, err := r.getBucket(tx)
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
	return r.Create(ctx, p) // Same for BoltDB
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

func (r *BoltRepository) List(ctx context.Context) ([]Profile, error) {
	var profiles []Profile
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getBucket(tx)
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
		b, err := r.getBucket(tx)
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
