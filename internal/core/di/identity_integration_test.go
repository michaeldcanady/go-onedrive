package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultContainer_IdentityWiring(t *testing.T) {
	// We need to provide a minimal environment for the container to initialize
	t.Setenv("HOME", t.TempDir())

	container, err := NewDefaultContainer()
	require.NoError(t, err)
	require.NotNil(t, container)

	t.Run("identity service is initialized", func(t *testing.T) {
		svc := container.Identity()
		assert.NotNil(t, svc)
	})

	t.Run("microsoft provider is registered", func(t *testing.T) {
		svc := container.Identity()
		providers := svc.ListProviders()
		assert.Contains(t, providers, "microsoft")
	})

	t.Run("microsoft authenticator is available", func(t *testing.T) {
		svc := container.Identity()
		auth, err := svc.GetAuthenticator("microsoft")
		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, "microsoft", auth.ProviderName())
	})
}
