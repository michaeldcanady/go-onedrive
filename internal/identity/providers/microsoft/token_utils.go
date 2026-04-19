package microsoft

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// extractIdentityFromToken extracts the identity from the access token.
// It looks for "preferred_username", "email", "upn", and finally "oid" claims.
func extractIdentityFromToken(tokenStr string) (string, error) {
	// A JWT has 3 parts separated by 2 dots.
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("token is not a JWT (found %d segments)", len(parts))
	}

	parser := jwt.NewParser()

	// We don't need to verify the signature here as we just got it from the provider.
	// We just want to extract the claims.
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims format")
	}

	// Try common identity claims
	if val, ok := claims["preferred_username"].(string); ok && val != "" {
		return val, nil
	}
	if val, ok := claims["email"].(string); ok && val != "" {
		return val, nil
	}
	if val, ok := claims["upn"].(string); ok && val != "" {
		return val, nil
	}
	if val, ok := claims["oid"].(string); ok && val != "" {
		return val, nil
	}

	return "", fmt.Errorf("no identity claim found in token")
}
