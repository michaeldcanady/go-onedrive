package microsoft

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestExtractIdentityFromToken(t *testing.T) {
	tests := []struct {
		name     string
		claims   jwt.MapClaims
		token    string
		expected string
		wantErr  bool
	}{
		{
			name: "preferred_username",
			claims: jwt.MapClaims{
				"preferred_username": "user@example.com",
			},
			expected: "user@example.com",
		},
		{
			name: "email",
			claims: jwt.MapClaims{
				"email": "user@example.com",
			},
			expected: "user@example.com",
		},
		{
			name: "upn",
			claims: jwt.MapClaims{
				"upn": "user@example.com",
			},
			expected: "user@example.com",
		},
		{
			name: "oid",
			claims: jwt.MapClaims{
				"oid": "12345",
			},
			expected: "12345",
		},
		{
			name: "priority_preferred_username",
			claims: jwt.MapClaims{
				"preferred_username": "user@example.com",
				"email":              "wrong@example.com",
			},
			expected: "user@example.com",
		},
		{
			name:    "no_claims",
			claims:  jwt.MapClaims{},
			wantErr: true,
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
				// We don't care about the signature for ParseUnverified
				tokenStr, _ = token.SignedString([]byte("secret"))
			}

			identity, err := extractIdentityFromToken(tokenStr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, identity)
			}
		})
	}
}
