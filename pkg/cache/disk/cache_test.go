package disk

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/pkg/cache"
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
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "test.cache")

	c, err := New[string, string](cachePath, stringSerde{}, stringSerde{})
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := "k1"
		val := "v1"

		err := c.Set(ctx, key, val)
		assert.NoError(t, err)

		got, err := c.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, val, got)
	})

	t.Run("Key Not Found", func(t *testing.T) {
		_, err := c.Get(ctx, "missing")
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("Update Entry", func(t *testing.T) {
		key := "k1"
		val2 := "v2"

		err := c.Set(ctx, key, val2)
		assert.NoError(t, err)

		got, err := c.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, val2, got)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "k-del"
		_ = c.Set(ctx, key, "val")

		err := c.Delete(ctx, key)
		assert.NoError(t, err)

		_, err = c.Get(ctx, key)
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("List", func(t *testing.T) {
		// Clear first to have clean state
		_ = c.Clear(ctx)

		_ = c.Set(ctx, "a", "1")
		_ = c.Set(ctx, "b", "2")

		results := make(map[string]string)
		err := c.List(ctx, func(k string, v string) error {
			results[k] = v
			return nil
		})

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "1", results["a"])
		assert.Equal(t, "2", results["b"])
	})

	t.Run("Persistence", func(t *testing.T) {
		key := "persist-key"
		val := "persist-val"
		_ = c.Set(ctx, key, val)

		// Create new instance pointing to same file
		c2, err := New[string, string](cachePath, stringSerde{}, stringSerde{})
		assert.NoError(t, err)

		got, err := c2.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, val, got)
	})

	t.Run("Clear", func(t *testing.T) {
		_ = c.Set(ctx, "x", "y")
		err := c.Clear(ctx)
		assert.NoError(t, err)

		_, err = c.Get(ctx, "x")
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)

		// File should be removed or empty
		_, err = os.Stat(cachePath)
		assert.True(t, os.IsNotExist(err))
	})
}
