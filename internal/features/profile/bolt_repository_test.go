package profile

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestBoltRepository(t *testing.T) {
	dbFile := "test_profiles.db"
	db, err := bolt.Open(dbFile, 0600, nil)
	assert.NoError(t, err)
	defer os.Remove(dbFile)
	defer db.Close()

	// Ensure bucket
	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte("profiles"))
		_, _ = tx.CreateBucketIfNotExists([]byte("settings"))
		return nil
	})

	repo := NewBoltRepository(db)
	ctx := context.Background()
	p := Profile{
		Name:      "test-profile",
		CreatedAt: time.Now(),
	}

	// Test Create
	err = repo.Create(ctx, p)
	assert.NoError(t, err)

	// Test Get
	retrieved, err := repo.Get(ctx, p.Name)
	assert.NoError(t, err)
	assert.Equal(t, p.Name, retrieved.Name)

	// Test Exists
	exists, err := repo.Exists(ctx, p.Name)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test List
	profiles, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, profiles, 1)

	// Test Settings
	err = repo.SetSetting(ctx, "key", "value")
	assert.NoError(t, err)
	val, err := repo.GetSetting(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	// Test Delete
	err = repo.Delete(ctx, p.Name)
	assert.NoError(t, err)
	_, err = repo.Get(ctx, p.Name)
	assert.Error(t, err)
}
