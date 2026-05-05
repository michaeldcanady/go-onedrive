package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetValue(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		setup   func() *Config
		verify  func(t *testing.T, c *Config)
		wantErr bool
	}{
		{
			name:  "set auth.provider",
			key:   "auth.provider",
			value: "google",
			setup: func() *Config { return &Config{} },
			verify: func(t *testing.T, c *Config) {
				assert.Equal(t, "google", c.Auth.Provider)
			},
		},
		{
			name:  "set logging.level",
			key:   "logging.level",
			value: "debug",
			setup: func() *Config { return &Config{} },
			verify: func(t *testing.T, c *Config) {
				assert.Equal(t, "debug", c.Logging.Level.String())
			},
		},
		{
			name:  "set editor.command",
			key:   "editor.command",
			value: "vim",
			setup: func() *Config { return &Config{} },
			verify: func(t *testing.T, c *Config) {
				assert.Equal(t, "vim", c.Editor.Command)
			},
		},
		{
			name:    "invalid key",
			key:     "invalid.key",
			value:   "val",
			setup:   func() *Config { return &Config{} },
			wantErr: true,
		},
		{
			name:    "missing subkey",
			key:     "auth",
			value:   "val",
			setup:   func() *Config { return &Config{} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			err := c.SetValue(tt.key, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.verify(t, c)
			}
		})
	}
}
