package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/abstractions"
)

type MetadataCache[K comparable, V any, M any] struct {
	abstractions.Cache[K, V]
	db          *sql.DB
	metadataSer abstractions.SerializerDeserializer[M]
}

func NewMetadataCache[K comparable, V any, M any](
	cache abstractions.Cache[K, V],
	path string,
	metadataSer abstractions.SerializerDeserializer[M],
) (*MetadataCache[K, V, M], error) {

	var db *sql.DB
	var err error
	if sqliteCache, ok := cache.(*Cache[K, V]); ok {
		db = sqliteCache.db
	}

	if db == nil {
		if db, err = sql.Open("sqlite3", path); err != nil {
			return nil, err
		}
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS metadata (
            key   BLOB PRIMARY KEY,
            value BLOB NOT NULL
        );
    `)
	if err != nil {
		return nil, err
	}

	return &MetadataCache[K, V, M]{
		db:          db,
		Cache:       cache,
		metadataSer: metadataSer,
	}, nil
}

func (c *MetadataCache[K, V, M]) GetMetadata(ctx context.Context, key K) (M, error) {

	var zero M

	serializedKey, err := c.KeySerializer().Serialize(key)
	if err != nil {
		return zero, err
	}

	row := c.db.QueryRowContext(ctx,
		`SELECT value FROM metadata WHERE key = ?`,
		serializedKey,
	)

	var raw []byte
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return zero, nil
		}
		return zero, err
	}

	md, err := c.metadataSer.Deserialize(raw)
	if err != nil {
		return zero, err
	}

	return md, nil
}

func (c *MetadataCache[K, V, M]) SetMetadata(ctx context.Context, key K, metadata M) error {

	serializedKey, err := c.KeySerializer().Serialize(key)
	if err != nil {
		return err
	}

	serializedValue, err := c.metadataSer.Serialize(metadata)
	if err != nil {
		return err
	}

	_, err = c.db.ExecContext(ctx,
		`INSERT INTO metadata(key, value)
         VALUES(?, ?)
         ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		serializedKey, serializedValue,
	)
	return err
}
