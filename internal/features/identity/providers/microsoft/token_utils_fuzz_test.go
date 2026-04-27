package microsoft

import (
	"testing"
)

func FuzzExtractFullIdentityFromToken(f *testing.F) {
	f.Add("header.payload.signature")
	f.Add("")
	f.Add("not-a-jwt")
	f.Add("a.b")
	f.Add("a.b.c.d")

	f.Fuzz(func(t *testing.T, tokenStr string) {
		_, _ = extractFullIdentityFromToken(tokenStr)
	})
}

func BenchmarkExtractFullIdentityFromToken(b *testing.B) {
	// A mock token that looks like a JWT
	// gitleaks:allow
	// nolint:gosec // G101 // not real credentials
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvaWQiOiJ1c2VyMSIsInByZWZlcnJlZF91c2VybmFtZSI6InVzZXIxQGV4YW1wbGUuY29tIiwibmFtZSI6IlVzZXIgMSJ9.signature"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = extractFullIdentityFromToken(token)
	}
}
