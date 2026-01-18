package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/abstractions"
)

type Cache[K comparable, V any] struct {
	db              *sql.DB
	keySerializer   abstractions.SerializerDeserializer[K]
	valueSerializer abstractions.SerializerDeserializer[V]
}

func New[K comparable, V any](
	path string,
	keySer abstractions.SerializerDeserializer[K],
	valueSer abstractions.SerializerDeserializer[V],
) (*Cache[K, V], error) {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS cache (
            key   BLOB PRIMARY KEY,
            value BLOB NOT NULL
        );
    `)
	if err != nil {
		return nil, err
	}

	return &Cache[K, V]{db, keySer, valueSer}, nil
}

func (c *Cache[K, V]) GetEntry(
	ctx context.Context,
	key K,
) (*abstractions.Entry[K, V], error) {

	serializedKey, err := c.keySerializer.Serialize(key)
	if err != nil {
		return nil, err
	}

	row := c.db.QueryRowContext(ctx,
		`SELECT value FROM cache WHERE key = ?`,
		serializedKey,
	)

	var raw []byte
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	value, err := c.valueSerializer.Deserialize(raw)
	if err != nil {
		return nil, err
	}

	return abstractions.NewEntry(key, value), nil
}

func (c *Cache[K, V]) SetEntry(
	ctx context.Context,
	entry *abstractions.Entry[K, V],
) error {

	serializedKey, err := c.keySerializer.Serialize(entry.GetKey())
	if err != nil {
		return err
	}

	serializedValue, err := c.valueSerializer.Serialize(entry.GetValue())
	if err != nil {
		return err
	}

	_, err = c.db.ExecContext(ctx,
		`INSERT INTO cache(key, value)
         VALUES(?, ?)
         ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		serializedKey, serializedValue,
	)
	return err
}

func (c *Cache[K, V]) Remove(key K) error {
	serializedKey, err := c.keySerializer.Serialize(key)
	if err != nil {
		return err
	}

	_, err = c.db.Exec(`DELETE FROM cache WHERE key = ?`, serializedKey)
	return err
}

func (c *Cache[K, V]) NewEntry(ctx context.Context, key K) (*abstractions.Entry[K, V], error) {
	entry, err := c.GetEntry(ctx, key)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		return nil, errors.New("key already exists")
	}

	var zero V
	return abstractions.NewEntry(key, zero), nil
}

func (c *Cache[K, V]) Clear(ctx context.Context) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM cache`)
	return err
}

func (c *Cache[K, V]) KeySerializer() abstractions.Serializer[K] {
	return c.keySerializer
}
