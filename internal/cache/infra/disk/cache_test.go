package infra

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/cache/infra/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stringSerializer struct{}

func (s *stringSerializer) Serialize(v string) ([]byte, error) {
	return []byte(v), nil
}

func (s *stringSerializer) Deserialize(b []byte) (string, error) {
	return string(b), nil
}

func TestCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "odc-disk-cache-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	cachePath := filepath.Join(tmpDir, "test.cache")
	ks := &stringSerializer{}
	vs := &stringSerializer{}

	c, err := New(cachePath, ks, vs)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		err := c.Set(ctx, "key1", "value1")
		assert.NoError(t, err)

		val, err := c.Get(ctx, "key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)
	})

	t.Run("Get Non-Existent", func(t *testing.T) {
		_, err := c.Get(ctx, "nonexistent")
		assert.ErrorIs(t, err, bolt.ErrKeyNotFound)
	})

	t.Run("Overwrite", func(t *testing.T) {
		err := c.Set(ctx, "key1", "value2")
		assert.NoError(t, err)

		val, err := c.Get(ctx, "key1")
		assert.NoError(t, err)
		assert.Equal(t, "value2", val)
	})

	t.Run("Delete", func(t *testing.T) {
		err := c.Delete(ctx, "key1")
		assert.NoError(t, err)

		_, err = c.Get(ctx, "key1")
		assert.ErrorIs(t, err, bolt.ErrKeyNotFound)
	})

	t.Run("Persistence", func(t *testing.T) {
		err := c.Set(ctx, "persist1", "val1")
		require.NoError(t, err)

		// Create new cache instance pointing to same file
		c2, err := New(cachePath, ks, vs)
		require.NoError(t, err)

		val, err := c2.Get(ctx, "persist1")
		assert.NoError(t, err)
		assert.Equal(t, "val1", val)
	})

	t.Run("List", func(t *testing.T) {
		err := c.Clear(ctx)
		require.NoError(t, err)

		entries := map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		}

		for k, v := range entries {
			err := c.Set(ctx, k, v)
			require.NoError(t, err)
		}

		got := make(map[string]string)
		err = c.List(ctx, func(k string, v string) error {
			got[k] = v
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, entries, got)
	})
}

func TestCache_Concurrency(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "odc-disk-cache-concurrency-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	cachePath := filepath.Join(tmpDir, "test.cache")
	ks := &stringSerializer{}
	vs := &stringSerializer{}

	c, err := New(cachePath, ks, vs)
	require.NoError(t, err)

	ctx := context.Background()
	const count = 100

	done := make(chan bool)
	for i := 0; i < count; i++ {
		go func(id int) {
			key := "key" // Same key to test locking
			val := "val"
			_ = c.Set(ctx, key, val)
			_, _ = c.Get(ctx, key)
			done <- true
		}(i)
	}

	for i := 0; i < count; i++ {
		<-done
	}
}
