package microsoft

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMicrosoftAuthenticator_ProviderName(t *testing.T) {
	a := NewMicrosoftAuthenticator()
	assert.Equal(t, "microsoft", a.ProviderName())
}

func TestMicrosoftAuthenticator_Logout(t *testing.T) {
	a := NewMicrosoftAuthenticator()
	ctx := context.Background()

	// Logout of a non-existent identity should not fail
	err := a.Logout(ctx, "non-existent")
	assert.NoError(t, err)

	// We can't easily verify the internal state of creds since it's private,
	// but we can ensure it doesn't panic and returns no error.
	a.creds["user1"] = new(mockTokenCredential)
	err = a.Logout(ctx, "user1")
	assert.NoError(t, err)
	
	a.mu.RLock()
	_, ok := a.creds["user1"]
	a.mu.RUnlock()
	assert.False(t, ok)
}
