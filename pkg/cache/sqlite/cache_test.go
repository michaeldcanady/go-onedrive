package sqlite

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type stringSerde struct{}

func (s stringSerde) Serialize(v string) ([]byte, error) {
	return []byte(v), nil
}

func (s stringSerde) Deserialize(b []byte) (string, error) {
	return string(b), nil
}

func TestCache(t *testing.T) {
	// Use in-memory sqlite
	c, err := New[string, string](":memory:", stringSerde{}, stringSerde{})
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := "k1"
		val := "v1"
		e := cache.NewEntry(key, val)

		err := c.SetEntry(ctx, e)
		assert.NoError(t, err)

		got, err := c.GetEntry(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, val, got.GetValue())
	})

	t.Run("Key Not Found", func(t *testing.T) {
		_, err := c.GetEntry(ctx, "missing")
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("Update Entry", func(t *testing.T) {
		key := "k1"
		val2 := "v2"
		e := cache.NewEntry(key, val2)

		err := c.SetEntry(ctx, e)
		assert.NoError(t, err)

		got, err := c.GetEntry(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, val2, got.GetValue())
	})

	t.Run("Remove", func(t *testing.T) {
		key := "k-rem"
		_ = c.SetEntry(ctx, cache.NewEntry(key, "val"))

		err := c.Remove(key)
		assert.NoError(t, err)

		_, err = c.GetEntry(ctx, key)
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("Clear", func(t *testing.T) {
		_ = c.SetEntry(ctx, cache.NewEntry("a", "1"))
		_ = c.SetEntry(ctx, cache.NewEntry("b", "2"))

		err := c.Clear(ctx)
		assert.NoError(t, err)

		_, err = c.GetEntry(ctx, "a")
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})
}
