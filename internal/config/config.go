package config

var _ Config = (*ConfigImpl)(nil)

type Config interface {
	GetAuthenticationConfig() AuthenticationConfig
	GetLoggingConfig() LoggingConfig
}

type ConfigImpl struct {
	Auth    *AuthenticationConfigImpl `mapstructure:"auth"`
	Logging *LoggingConfigImpl        `mapstructure:"logging"`
}

func (c *ConfigImpl) GetAuthenticationConfig() AuthenticationConfig {
	return c.Auth
}

func (c *ConfigImpl) GetLoggingConfig() LoggingConfig {
	return c.Logging
}
