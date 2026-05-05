package identity

import (
	"context"
	"testing"
)

func FuzzParseAuthMethod(f *testing.F) {
	f.Add("interactive")
	f.Add("device-code")
	f.Add("client-secret")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		_ = ParseAuthMethod(input)
	})
}

func FuzzRegistry_Login(f *testing.F) {
	ctx := context.Background()
	reg := NewRegistry(nil, nil)

	f.Add("microsoft")
	f.Add("")

	f.Fuzz(func(t *testing.T, provider string) {
		// Just ensure it doesn't panic on random provider strings
		_, _ = reg.Login(ctx, provider, LoginOptions{})
	})
}
