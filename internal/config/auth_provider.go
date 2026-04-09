package config

import "encoding/json"

type AuthProvider int8

const (
	AuthProviderUnknown AuthProvider = iota
	AuthProviderMicrosoft
)

func (p AuthProvider) String() string {
	switch p {
	case AuthProviderMicrosoft:
		return "microsoft"
	default:
		return "unknown"
	}
}

func ParseAuthProvider(s string) AuthProvider {
	switch s {
	case "microsoft":
		return AuthProviderMicrosoft
	default:
		return AuthProviderUnknown
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (p AuthProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *AuthProvider) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*p = ParseAuthProvider(s)
	return nil
}

// MarshalText implements encoding.TextMarshaler for serializing AuthProvider as a string.
func (p AuthProvider) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for deserializing AuthProvider from a string.
func (p *AuthProvider) UnmarshalText(text []byte) error {
	*p = ParseAuthProvider(string(text))
	return nil
}
