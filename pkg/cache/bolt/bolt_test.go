package bolt

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	bucketName := "test-bucket"

	s, err := NewStore(dbPath, bucketName)
	assert.NoError(t, err)
	defer s.Close()

	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := []byte("key1")
		value := []byte("value1")

		err := s.Set(ctx, key, value)
		assert.NoError(t, err)

		got, err := s.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, got)
	})

	t.Run("Get Non-existent Key", func(t *testing.T) {
		_, err := s.Get(ctx, []byte("non-existent"))
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("Delete", func(t *testing.T) {
		key := []byte("to-delete")
		_ = s.Set(ctx, key, []byte("val"))

		err := s.Delete(ctx, key)
		assert.NoError(t, err)

		_, err = s.Get(ctx, key)
		assert.ErrorIs(t, err, cache.ErrKeyNotFound)
	})

	t.Run("List", func(t *testing.T) {
		// Clear or use fresh bucket for predictable List
		bucket2 := "list-bucket"
		s2, err := NewSiblingStore(s, bucket2)
		assert.NoError(t, err)

		_ = s2.Set(ctx, []byte("a"), []byte("1"))
		_ = s2.Set(ctx, []byte("b"), []byte("2"))

		keys, values, err := s2.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, keys, 2)
		assert.Equal(t, [][]byte{[]byte("a"), []byte("b")}, keys)
		assert.Equal(t, [][]byte{[]byte("1"), []byte("2")}, values)
	})
}
