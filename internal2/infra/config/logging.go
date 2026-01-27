package config

type LoggingConfig interface {
	// GetLevel returns the logging level.
	GetLevel() string
}

type LoggingConfigImpl struct {
	Level string `mapstructure:"level"`
}

// GetLevel returns the logging level.
func (l *LoggingConfigImpl) GetLevel() string {
	return l.Level
}
