package identity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/errors"
	bolt "go.etcd.io/bbolt"
)

// BoltRepository implements IdentityRepository using BoltDB.
// BoltRepository implements IdentityRepository using BoltDB.
type BoltRepository struct {
	db *bolt.DB
}

// NewBoltRepository creates a new instance of BoltRepository.
const (
	tokensBucketName = "tokens"
)

// NewBoltRepository creates a new instance of BoltRepository.
func NewBoltRepository(db *bolt.DB) *BoltRepository {
	return &BoltRepository{db: db}
}

// Initialize ensures the DB schema (tokens bucket) exists.
func (r *BoltRepository) Initialize() error {
	return r.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(tokensBucketName))
		return err
	})
}

// getProviderBucket returns the bucket for a specific provider.
func (r *BoltRepository) getProviderBucket(tx *bolt.Tx, provider string) (*bolt.Bucket, error) {
	root := tx.Bucket([]byte(tokensBucketName))
	if root == nil {
		if tx.Writable() {
			var err error
			root, err = tx.CreateBucketIfNotExists([]byte(tokensBucketName))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("tokens bucket missing")
		}
	}
	if tx.Writable() {
		return root.CreateBucketIfNotExists([]byte(provider))
	}
	return root.Bucket([]byte(provider)), nil
}

func (r *BoltRepository) Get(ctx context.Context, provider, identityID string) (AccessToken, error) {
	var token AccessToken
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getProviderBucket(tx, provider)
		if err != nil {
			return err
		}
		if b == nil {
			return fmt.Errorf("%w: provider bucket %s missing", errors.ErrNotFound, provider)
		}
		data := b.Get([]byte(identityID))
		if data == nil {
			return fmt.Errorf("%w for identity %s in provider %s", errors.ErrNotFound, identityID, provider)
		}
		return json.Unmarshal(data, &token)
	})
	return token, err
}

func (r *BoltRepository) Save(ctx context.Context, provider string, token AccessToken) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := r.getProviderBucket(tx, provider)
		if err != nil {
			return err
		}
		// nolint:gosec // G117 // allowed
		data, err := json.Marshal(token)
		if err != nil {
			return fmt.Errorf("failed to marshal token: %w", err)
		}
		return b.Put([]byte(token.AccountID), data)
	})
}

func (r *BoltRepository) Delete(ctx context.Context, provider, AccountID string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := r.getProviderBucket(tx, provider)
		if err != nil {
			return err
		}
		if b == nil {
			return nil
		}
		return b.Delete([]byte(AccountID))
	})
}

func (r *BoltRepository) List(ctx context.Context, provider string) ([]string, error) {
	var ids []string
	err := r.db.View(func(tx *bolt.Tx) error {
		b, err := r.getProviderBucket(tx, provider)
		if err != nil {
			return err
		}
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			if v != nil {
				ids = append(ids, string(k))
			}
			return nil
		})
	})
	return ids, err
}
