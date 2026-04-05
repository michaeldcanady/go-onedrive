package config

// Config represents the application's top-level configuration structure.
type Config struct {
	// Auth contains settings related to authentication and provider identity.
	Auth AuthenticationConfig `json:"auth" yaml:"auth"`
	// Logging contains settings related to logging behavior and output.
	Logging LoggingConfig `json:"logging" yaml:"logging"`
}
