package identity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMethod_String(t *testing.T) {
	tests := []struct {
		name string
		m    AuthMethod
		want string
	}{
		{"interactive", AuthMethodInteractiveBrowser, "interactive"},
		{"device-code", AuthMethodDeviceCode, "device-code"},
		{"client-secret", AuthMethodClientSecret, "client-secret"},
		{"environment", AuthMethodEnvironment, "environment"},
		{"unknown", AuthMethodUnknown, "unknown"},
		{"invalid", AuthMethod(99), "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.String())
		})
	}
}

func TestParseAuthMethod(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  AuthMethod
	}{
		{"interactive", "interactive", AuthMethodInteractiveBrowser},
		{"browser", "browser", AuthMethodInteractiveBrowser},
		{"device-code", "device-code", AuthMethodDeviceCode},
		{"device", "device", AuthMethodDeviceCode},
		{"client-secret", "client-secret", AuthMethodClientSecret},
		{"secret", "secret", AuthMethodClientSecret},
		{"environment", "environment", AuthMethodEnvironment},
		{"env", "env", AuthMethodEnvironment},
		{"unknown", "unknown", AuthMethodUnknown},
		{"invalid", "invalid", AuthMethodUnknown},
		{"case insensitive", "INTERACTIVE", AuthMethodInteractiveBrowser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseAuthMethod(tt.input))
		})
	}
}

func TestAuthMethod_JSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		m := AuthMethodInteractiveBrowser
		data, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `"interactive"`, string(data))
	})

	t.Run("unmarshal", func(t *testing.T) {
		var m AuthMethod
		err := json.Unmarshal([]byte(`"device-code"`), &m)
		assert.NoError(t, err)
		assert.Equal(t, AuthMethodDeviceCode, m)
	})
}
