package identity

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	identitiesBucket = []byte("identities")
	tokensBucket     = []byte("tokens")
)

type boltRepository struct {
	db *bbolt.DB
}

// NewBoltRepository creates a new bbolt-based identity repository.
func NewBoltRepository(db *bbolt.DB) (Repository, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(identitiesBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(tokensBucket)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize identity buckets: %w", err)
	}

	return &boltRepository{db: db}, nil
}

func (r *boltRepository) SaveIdentity(i *Identity) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(identitiesBucket)
		data, err := json.Marshal(i)
		if err != nil {
			return err
		}
		return b.Put([]byte(i.ID), data)
	})
}

func (r *boltRepository) GetIdentity(id string) (*Identity, error) {
	var i Identity
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(identitiesBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &i)
	})
	if i.ID == "" {
		return nil, nil
	}
	return &i, err
}

func (r *boltRepository) ListIdentities() ([]*Identity, error) {
	var identities []*Identity
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(identitiesBucket)
		return b.ForEach(func(k, v []byte) error {
			var i Identity
			if err := json.Unmarshal(v, &i); err != nil {
				return err
			}
			identities = append(identities, &i)
			return nil
		})
	})
	return identities, err
}

func (r *boltRepository) DeleteIdentity(id string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(identitiesBucket)
		return b.Delete([]byte(id))
	})
}

func (r *boltRepository) SaveToken(provider string, identityID string, t *Token) error {
	key := fmt.Sprintf("%s:%s", provider, identityID)
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tokensBucket)
		data, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), data)
	})
}

func (r *boltRepository) GetToken(provider string, identityID string) (*Token, error) {
	key := fmt.Sprintf("%s:%s", provider, identityID)
	var t Token
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tokensBucket)
		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &t)
	})
	if t.AccessToken == "" {
		return nil, nil
	}
	return &t, err
}

func (r *boltRepository) DeleteToken(provider string, identityID string) error {
	key := fmt.Sprintf("%s:%s", provider, identityID)
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tokensBucket)
		return b.Delete([]byte(key))
	})
}
