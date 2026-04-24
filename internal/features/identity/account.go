package identity

// Account represents a structured user identity from an identity provider.
type Account struct {
	// ID is the unique identifier for the identity (e.g., a UUID or UPN).
	ID string `json:"id" yaml:"id"`
	// DisplayName is the human-readable name of the user.
	DisplayName string `json:"display_name" yaml:"display_name"`
	// Email is the user's email address.
	Email string `json:"email" yaml:"email"`
	// Provider is the name of the identity provider (e.g., "microsoft").
	Provider string `json:"provider" yaml:"provider"`
	// AvatarURL is an optional URL to the user's profile picture.
	AvatarURL string `json:"avatar_url,omitempty" yaml:"avatar_url,omitempty"`
	// Metadata contains provider-specific additional information.
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}
