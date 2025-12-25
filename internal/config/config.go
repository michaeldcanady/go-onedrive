package config

type Config interface {
	GetAuthenticationConfig() AuthenticationConfig
}
