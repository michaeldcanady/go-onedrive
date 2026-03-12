package config

// AuthenticationConfig defines the settings required for provider-specific authentication.
type AuthenticationConfig struct {
	// Provider identifies the authentication provider (e.g., "microsoft", "google").
	Provider     string `json:"provider" yaml:"provider" mapstructure:"provider"`
	// ClientID is the unique identifier for the application registered with the provider.
	ClientID     string `json:"client_id" yaml:"client_id" mapstructure:"client_id"`
	// TenantID is the organizational identifier (specific to Microsoft Graph).
	TenantID     string `json:"tenant_id" yaml:"tenant_id" mapstructure:"tenant_id"`
	// ClientSecret is the secret used for Service Principal authentication.
	ClientSecret string `json:"client_secret" yaml:"client_secret" mapstructure:"client_secret"`
	// Method specifies the authentication mechanism to use.
	Method       string `json:"method" yaml:"method" mapstructure:"method"`
	// RedirectURI is the endpoint where authentication responses are received.
	RedirectURI  string `json:"redirect_uri" yaml:"redirect_uri" mapstructure:"redirect_uri"`
}
