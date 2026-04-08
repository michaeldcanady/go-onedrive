package memory

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Run("GetEntry", func(t *testing.T) {
		c := New[string, string]()
		ctx := context.Background()
		key := "test-key"
		value := "test-value"

		// Key not found
		entry, err := c.GetEntry(ctx, key)
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
		assert.Nil(t, entry)

		// Set and Get
		e := cache.NewEntry(key, value)
		err = c.SetEntry(ctx, e)
		assert.NoError(t, err)

		entry, err = c.GetEntry(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, entry.GetValue())
	})

	t.Run("NewEntry", func(t *testing.T) {
		c := New[string, int]()
		ctx := context.Background()
		key := "new-key"

		entry, err := c.NewEntry(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, key, entry.GetKey())
		assert.Equal(t, 0, entry.GetValue())
	})

	t.Run("Remove", func(t *testing.T) {
		c := New[string, string]()
		ctx := context.Background()
		key := "remove-key"

		e := cache.NewEntry(key, "val")
		_ = c.SetEntry(ctx, e)

		err := c.Remove(key)
		assert.NoError(t, err)

		_, err = c.GetEntry(ctx, key)
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("Clear", func(t *testing.T) {
		c := New[string, string]()
		ctx := context.Background()

		_ = c.SetEntry(ctx, cache.NewEntry("k1", "v1"))
		_ = c.SetEntry(ctx, cache.NewEntry("k2", "v2"))

		err := c.Clear(ctx)
		assert.NoError(t, err)

		_, err = c.GetEntry(ctx, "k1")
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("ContextCanceled", func(t *testing.T) {
		c := New[string, string]()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := c.GetEntry(ctx, "any")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")

		err = c.SetEntry(ctx, cache.NewEntry("k", "v"))
		assert.Error(t, err)

		err = c.Clear(ctx)
		assert.Error(t, err)
	})
}
