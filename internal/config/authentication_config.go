package config

import "github.com/michaeldcanady/go-onedrive/internal/identity/shared"

// AuthenticationConfig defines the settings required for provider-specific authentication.
type AuthenticationConfig struct {
	// Provider identifies the authentication provider (e.g., "microsoft", "google").
	Provider AuthProvider `json:"provider,omitempty" yaml:"provider,omitempty" mapstructure:"provider"`
	// ClientID is the unique identifier for the application registered with the provider.
	ClientID string `json:"client_id,omitempty" yaml:"client_id,omitempty" mapstructure:"client_id"`
	// TenantID is the organizational identifier (specific to Microsoft Graph).
	TenantID string `json:"tenant_id,omitempty" yaml:"tenant_id,omitempty" mapstructure:"tenant_id"`
	// ClientSecret is the secret used for Service Principal authentication.
	ClientSecret string `json:"client_secret,omitempty" yaml:"client_secret,omitempty" mapstructure:"client_secret"`
	// Method specifies the authentication mechanism to use.
	Method shared.AuthMethod `json:"method,omitempty" yaml:"method,omitempty" mapstructure:"method"`
	// RedirectURI is the endpoint where authentication responses are received.
	RedirectURI string `json:"redirect_uri,omitempty" yaml:"redirect_uri,omitempty" mapstructure:"redirect_uri"`
}
