package bolt

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	bbolt "go.etcd.io/bbolt"
)

type Store struct {
	db     *bbolt.DB
	bucket []byte
}

func NewStore(path string, bucket string) (*Store, error) {
	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	}); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Store{
		db:     db,
		bucket: []byte(bucket),
	}, nil
}

func NewSiblingStore(store *Store, bucket string) (*Store, error) {
	if err := store.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	}); err != nil {
		return nil, err
	}

	return &Store{
		db:     store.db,
		bucket: []byte(bucket),
	}, nil
}

func NewStoreWithDB(db *bbolt.DB, bucket string) (*Store, error) {
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	}); err != nil {
		return nil, err
	}

	return &Store{
		db:     db,
		bucket: []byte(bucket),
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	var out []byte

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bucket)
		if b == nil {
			return core.ErrKeyNotFound
		}

		v := b.Get(key)
		if v == nil {
			return core.ErrKeyNotFound
		}

		out = append([]byte(nil), v...)
		return nil
	})

	return out, err
}

func (s *Store) Set(ctx context.Context, key, value []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bucket)
		if b == nil {
			return core.ErrKeyNotFound
		}
		return b.Put(key, value)
	})
}

func (s *Store) Delete(ctx context.Context, key []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(s.bucket)
		if b == nil {
			return core.ErrKeyNotFound
		}
		return b.Delete(key)
	})
}
