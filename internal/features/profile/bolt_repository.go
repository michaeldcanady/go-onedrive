package profile

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	profilesBucket = []byte("profiles")
	metaBucket     = []byte("meta")
	currentProfile = []byte("current_profile")
)

type boltRepository struct {
	db *bbolt.DB
}

func initializeRepository(db *bbolt.DB) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(profilesBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(metaBucket)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to initialize profile buckets: %w", err)
	}
	return nil
}

// NewBoltRepository creates a new bbolt-based profile repository.
func NewBoltRepository(db *bbolt.DB) (Repository, error) {
	if err := initializeRepository(db); err != nil {
		return nil, err
	}
	return &boltRepository{db: db}, nil
}

func (r *boltRepository) Create(p *Profile) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(profilesBucket)
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(p.Name), data)
	})
}

func (r *boltRepository) List() ([]*Profile, error) {
	var profiles []*Profile
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(profilesBucket)
		return b.ForEach(func(k, v []byte) error {
			var p Profile
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			profiles = append(profiles, &p)
			return nil
		})
	})
	return profiles, err
}

func (r *boltRepository) Delete(name string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(profilesBucket)
		return b.Delete([]byte(name))
	})
}

func (r *boltRepository) GetCurrent() (string, error) {
	var name string
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(metaBucket)
		v := b.Get(currentProfile)
		if v != nil {
			name = string(v)
		}
		return nil
	})
	return name, err
}

func (r *boltRepository) SetCurrent(name string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(metaBucket)
		return b.Put(currentProfile, []byte(name))
	})
}
