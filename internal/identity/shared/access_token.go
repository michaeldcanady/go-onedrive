package shared

import "time"

// AccessToken represents an authentication token and its associated metadata.
type AccessToken struct {
	// IdentityID is the unique identifier for the user or account (e.g., "user@outlook.com").
	IdentityID string `json:"identity_id" yaml:"identity_id"`
	// Token is the raw access token string.
	Token string `json:"token" yaml:"token"`
	// RefreshToken is the token used to obtain a new access token.
	RefreshToken string `json:"refresh_token" yaml:"refresh_token"`
	// ExpiresAt is the timestamp when the access token becomes invalid.
	ExpiresAt time.Time `json:"expires_at" yaml:"expires_at"`
	// Scopes are the permissions granted by this token.
	Scopes []string `json:"scopes" yaml:"scopes"`
	// ProviderSpecific contains any extra data provided by the identity provider.
	ProviderSpecific map[string]any `json:"provider_specific" yaml:"provider_specific"`
}
