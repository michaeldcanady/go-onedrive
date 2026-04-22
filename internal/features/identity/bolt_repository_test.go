package identity

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestBoltRepository(t *testing.T) {
	dbFile := "test_identity.db"
	db, err := bolt.Open(dbFile, 0600, nil)
	assert.NoError(t, err)
	defer os.Remove(dbFile)
	defer db.Close()

	repo := NewBoltRepository(db)
	ctx := context.Background()
	provider := "microsoft"
	token := AccessToken{
		AccountID: "user@test.com",
		Token:     "secret-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Test Save
	err = repo.Save(ctx, provider, token)
	assert.NoError(t, err)

	// Test Get
	retrieved, err := repo.Get(ctx, provider, token.AccountID)
	assert.NoError(t, err)
	assert.Equal(t, token.Token, retrieved.Token)

	// Test List
	ids, err := repo.List(ctx, provider)
	assert.NoError(t, err)
	assert.Contains(t, ids, token.AccountID)

	// Test Delete
	err = repo.Delete(ctx, provider, token.AccountID)
	assert.NoError(t, err)

	_, err = repo.Get(ctx, provider, token.AccountID)
	assert.Error(t, err)
}
