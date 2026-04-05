package config

import (
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

type LoggingConfig struct {
	// Level specifies the minimum severity level of log messages.
	Level logger.Level `json:"level,omitempty" yaml:"level,omitempty" mapstructure:"level"`
	// Output specifies the destination for log messages (e.g., "stdout", "stderr", or a file path).
	Output string `json:"output,omitempty" yaml:"output,omitempty" mapstructure:"output"`
	// Format specifies the log message format (e.g., "json" or "text").
	Format string `json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format"`
}
