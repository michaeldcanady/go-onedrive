package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAvailableKeys(t *testing.T) {
	keys := GetAvailableKeys()
	assert.NotEmpty(t, keys)
	assert.Contains(t, keys, "auth.client_id")
	assert.Contains(t, keys, "auth.provider")
	assert.Contains(t, keys, "logging.level")
	assert.Contains(t, keys, "logging.output")
	assert.Contains(t, keys, "logging.format")
}

func TestGetAllowedValues(t *testing.T) {
	tests := []struct {
		key      string
		expected []string
	}{
		{
			key:      "auth.provider",
			expected: []string{"microsoft"},
		},
		{
			key:      "logging.level",
			expected: []string{"debug", "error", "fatal", "info", "warn"},
		},
		{
			key:      "logging.format",
			expected: []string{"json", "text"},
		},
		{
			key:      "logging.output",
			expected: []string{"stderr", "stdout"},
		},
		{
			key:      "invalid.key",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			values := GetAllowedValues(tt.key)
			assert.Equal(t, tt.expected, values)
		})
	}
}
