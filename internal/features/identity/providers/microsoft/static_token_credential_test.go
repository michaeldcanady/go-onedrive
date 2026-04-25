package microsoft

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/stretchr/testify/assert"
)

func TestStaticTokenCredential_GetToken(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		token     identity.AccessToken
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid token",
			token: identity.AccessToken{
				Token:     "valid-token",
				ExpiresAt: now.Add(1 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "expired token",
			token: identity.AccessToken{
				Token:     "expired-token",
				ExpiresAt: now.Add(-1 * time.Hour),
			},
			wantErr: true,
			errMsg:  "cached access token is expired",
		},
		{
			name: "token expiring soon",
			token: identity.AccessToken{
				Token:     "expiring-soon-token",
				ExpiresAt: now.Add(4 * time.Minute),
			},
			wantErr: true,
			errMsg:  "cached access token is expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := NewStaticTokenCredential(tt.token)
			got, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.token.Token, got.Token)
				assert.Equal(t, tt.token.ExpiresAt, got.ExpiresOn)
			}
		})
	}
}
