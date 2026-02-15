package auth

import "time"

type AccessToken struct {
	// Token is the access token
	Token string `json:"token"`
	// ExpiresOn indicates when the token expires
	ExpiresOn time.Time `json:"expires_on"`
	// RefreshOn is a suggested time to refresh the token.
	// Clients should ignore this value when it's zero.
	RefreshOn time.Time `json:"refresh_on"`
}
