package microsoft

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
)

// extractFullIdentityFromToken extracts rich identity information from the access token.
func extractFullIdentityFromToken(tokenStr string) (identity.Account, error) {
	var ident identity.Account

	// A JWT has 3 parts separated by 2 dots.
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return ident, fmt.Errorf("token is not a JWT (found %d segments)", len(parts))
	}

	parser := jwt.NewParser()

	// We don't need to verify the signature here as we just got it from the provider.
	// We just want to extract the claims.
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return ident, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ident, fmt.Errorf("invalid claims format")
	}

	// Extract Provider-specific ID (oid)
	if val, ok := claims["oid"].(string); ok && val != "" {
		ident.ID = val
	}

	// Extract Email / Username
	if val, ok := claims["preferred_username"].(string); ok && val != "" {
		ident.Email = val
	} else if val, ok := claims["email"].(string); ok && val != "" {
		ident.Email = val
	} else if val, ok := claims["upn"].(string); ok && val != "" {
		ident.Email = val
	}

	// Extract Display Name
	if val, ok := claims["name"].(string); ok && val != "" {
		ident.DisplayName = val
	} else if val, ok := claims["name"].(string); ok && val != "" {
		ident.DisplayName = val
	}

	// Fallbacks
	if ident.Email == "" && ident.ID != "" {
		ident.Email = ident.ID
	}
	if ident.DisplayName == "" {
		ident.DisplayName = ident.Email
	}

	ident.Provider = "microsoft"

	// Check if we found anything useful
	if ident.Email == "" && ident.ID == "" {
		return ident, fmt.Errorf("no identity claims found in token")
	}

	return ident, nil
}
