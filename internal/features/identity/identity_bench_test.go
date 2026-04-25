package identity

import (
	"testing"
)

func BenchmarkParseAuthMethod(b *testing.B) {
	inputs := []string{"interactive", "browser", "device-code", "device", "client-secret", "secret", "environment", "env", "unknown", "invalid"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseAuthMethod(inputs[i%len(inputs)])
	}
}

func BenchmarkRegistry_ListProviders(b *testing.B) {
	reg := NewRegistry(nil, nil)
	reg.RegisterAuthenticator("p1", nil)
	reg.RegisterAuthenticator("p2", nil)
	reg.RegisterAuthorizer("p3", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reg.ListProviders()
	}
}
