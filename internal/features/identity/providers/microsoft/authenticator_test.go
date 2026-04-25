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
	tests := []struct {
		name       string
		identityID string
		setup      func(a *MicrosoftAuthenticator)
		wantErr    bool
	}{
		{
			name:       "logout non-existent identity",
			identityID: "non-existent",
			setup:      func(a *MicrosoftAuthenticator) {},
			wantErr:    false,
		},
		{
			name:       "logout existing identity",
			identityID: "user1",
			setup: func(a *MicrosoftAuthenticator) {
				a.mu.Lock()
				a.creds["user1"] = new(mockTokenCredential)
				a.mu.Unlock()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewMicrosoftAuthenticator()
			tt.setup(a)

			err := a.Logout(context.Background(), tt.identityID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				a.mu.RLock()
				_, ok := a.creds[tt.identityID]
				a.mu.RUnlock()
				assert.False(t, ok)
			}
		})
	}
}
