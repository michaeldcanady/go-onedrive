package config

import "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"

type Configuration struct {
	Auth AuthenticationConfig `json:"auth" yaml:"auth"`
}

type AuthenticationConfig struct {
	Type         auth.Method `json:"type" yaml:"type" mapstructure:"type"`
	ClientID     string      `json:"client_id" yaml:"client_id" mapstructure:"client_id"`
	TenantID     string      `json:"tenant_id" yaml:"tenant_id" mapstructure:"tenant_id"`
	RedirectURI  string      `json:"redirect_uri" yaml:"redirect_uri" mapstructure:"redirect_uri"`
	ProfileCache string      `json:"profile_cache" yaml:"profile_cache" mapstructure:"profile_cache"`
}

type LoggingConfig struct {
	Level string `json:"level" yaml:"level" mapstructure:"level"`
}
