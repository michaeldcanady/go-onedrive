package identity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/errors"
	bolt "go.etcd.io/bbolt"
)

type BoltRepository struct {
	db *bolt.DB
}

func NewBoltRepository(db *bolt.DB) *BoltRepository {
	return &BoltRepository{db: db}
}

func (r *BoltRepository) getBucket(tx *bolt.Tx, provider string) (*bolt.Bucket, error) {
	root, err := tx.CreateBucketIfNotExists([]byte("tokens"))
	if err != nil {
		return nil, err
	}
	return root.CreateBucketIfNotExists([]byte(provider))
}

func (r *BoltRepository) Get(ctx context.Context, provider, identityID string) (AccessToken, error) {
	var token AccessToken
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tokens"))
		if b == nil {
			return fmt.Errorf("%w for identity %s: tokens bucket missing", errors.ErrNotFound, identityID)
		}
		pb := b.Bucket([]byte(provider))
		if pb == nil {
			return fmt.Errorf("%w for identity %s: provider bucket %s missing", errors.ErrNotFound, identityID, provider)
		}
		data := pb.Get([]byte(identityID))
		if data == nil {
			return fmt.Errorf("%w for identity %s", errors.ErrNotFound, identityID)
		}
		return json.Unmarshal(data, &token)
	})
	return token, err
}

func (r *BoltRepository) Save(ctx context.Context, provider string, token AccessToken) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		pb, err := r.getBucket(tx, provider)
		if err != nil {
			return err
		}
		data, err := json.Marshal(token)
		if err != nil {
			return err
		}
		return pb.Put([]byte(token.AccountID), data)
	})
}

func (r *BoltRepository) Delete(ctx context.Context, provider, AccountID string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tokens"))
		if b == nil {
			return nil
		}
		pb := b.Bucket([]byte(provider))
		if pb == nil {
			return nil
		}
		return pb.Delete([]byte(AccountID))
	})
}

func (r *BoltRepository) List(ctx context.Context, provider string) ([]string, error) {
	var ids []string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tokens"))
		if b == nil {
			return nil
		}
		pb := b.Bucket([]byte(provider))
		if pb == nil {
			return nil
		}
		return pb.ForEach(func(k, v []byte) error {
			if v != nil {
				ids = append(ids, string(k))
			}
			return nil
		})
	})
	return ids, err
}
