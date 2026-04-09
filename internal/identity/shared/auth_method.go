package shared

import (
	"encoding/json"
	"strings"
)

// AuthMethod represents the mechanism used for authentication.
type AuthMethod int

const (
	// AuthMethodUnknown represents an unspecified authentication method.
	AuthMethodUnknown AuthMethod = iota
	// AuthMethodInteractiveBrowser uses the system's default web browser for login.
	AuthMethodInteractiveBrowser
	// AuthMethodDeviceCode provides a code for the user to enter on a separate device.
	AuthMethodDeviceCode
	// AuthMethodClientSecret uses a client ID and secret (Service Principal).
	AuthMethodClientSecret
	// AuthMethodEnvironment uses environment variables for credentials.
	AuthMethodEnvironment
)

// String returns the string representation of the authentication method.
func (m AuthMethod) String() string {
	switch m {
	case AuthMethodInteractiveBrowser:
		return "interactive"
	case AuthMethodDeviceCode:
		return "device-code"
	case AuthMethodClientSecret:
		return "client-secret"
	case AuthMethodEnvironment:
		return "environment"
	default:
		return "unknown"
	}
}

// ParseAuthMethod converts a string to its corresponding AuthMethod.
func ParseAuthMethod(s string) AuthMethod {
	switch strings.ToLower(s) {
	case "interactive", "browser":
		return AuthMethodInteractiveBrowser
	case "device-code", "device":
		return AuthMethodDeviceCode
	case "client-secret", "secret":
		return AuthMethodClientSecret
	case "environment", "env":
		return AuthMethodEnvironment
	default:
		return AuthMethodUnknown
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (m AuthMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *AuthMethod) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*m = ParseAuthMethod(s)
	return nil
}

// MarshalText implements encoding.TextMarshaler for serializing AuthMethod as a string.
func (m AuthMethod) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for deserializing AuthMethod from a string.
func (m *AuthMethod) UnmarshalText(text []byte) error {
	*m = ParseAuthMethod(string(text))
	return nil
}
