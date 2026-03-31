package config

// Config represents the application's top-level configuration structure.
type Config struct {
	// Auth contains settings related to authentication and provider identity.
	Auth AuthenticationConfig `json:"auth" yaml:"auth"`
}
