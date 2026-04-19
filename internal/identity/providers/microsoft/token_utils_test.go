package microsoft

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/stretchr/testify/assert"
)

func TestExtractFullIdentityFromToken(t *testing.T) {
	tests := []struct {
		name     string
		claims   jwt.MapClaims
		token    string
		expected identity.Account
		wantErr  bool
	}{
		{
			name: "full_claims",
			claims: jwt.MapClaims{
				"preferred_username": "user@example.com",
				"name":               "Test User",
				"oid":                "unique-oid",
			},
			expected: identity.Account{
				ID:          "unique-oid",
				Email:       "user@example.com",
				DisplayName: "Test User",
				Provider:    "microsoft",
			},
		},
		{
			name: "preferred_username_only",
			claims: jwt.MapClaims{
				"preferred_username": "user@example.com",
			},
			expected: identity.Account{
				ID:          "",
				Email:       "user@example.com",
				DisplayName: "user@example.com",
				Provider:    "microsoft",
			},
		},
		{
			name: "oid_only",
			claims: jwt.MapClaims{
				"oid": "unique-oid",
			},
			expected: identity.Account{
				ID:          "unique-oid",
				Email:       "unique-oid",
				DisplayName: "unique-oid",
				Provider:    "microsoft",
			},
		},
		{
			name:    "malformed_token",
			token:   "not-a-jwt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr := tt.token
			if tokenStr == "" && tt.claims != nil {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)
				tokenStr, _ = token.SignedString([]byte("secret"))
			}

			identity, err := extractFullIdentityFromToken(tokenStr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, identity.ID)
				assert.Equal(t, tt.expected.Email, identity.Email)
				assert.Equal(t, tt.expected.DisplayName, identity.DisplayName)
				assert.Equal(t, tt.expected.Provider, identity.Provider)
			}
		})
	}
}
