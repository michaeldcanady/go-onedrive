package di

import (
	"os"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/environment"
	"github.com/stretchr/testify/assert"
)

func setupTestEnv(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv(environment.EnvConfigDir, tmpDir)
	os.Setenv(environment.EnvCacheDir, tmpDir)
	os.Setenv(environment.EnvDataDir, tmpDir)
	os.Setenv(environment.EnvLogDir, tmpDir)
	os.Setenv(environment.EnvStateDir, tmpDir)
}

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	assert.NotNil(t, c)
}

func TestContainer_Methods(t *testing.T) {
	setupTestEnv(t)
	c := NewContainer()

	tests := []struct {
		name string
		fn   func(*Container) interface{}
	}{
		// Public Services
		{"EnvironmentService", func(c *Container) interface{} { return c.EnvironmentService() }},
		{"Logger", func(c *Container) interface{} { return c.Logger() }},
		{"Config", func(c *Container) interface{} { return c.Config() }},
		{"State", func(c *Container) interface{} { return c.State() }},
		{"Cache", func(c *Container) interface{} { return c.Cache() }},
		{"Auth", func(c *Container) interface{} { return c.Auth() }},
		{"Account", func(c *Container) interface{} { return c.Account() }},
		{"Profile", func(c *Container) interface{} { return c.Profile() }},
		{"Drive", func(c *Container) interface{} { return c.Drive() }},
		{"FS", func(c *Container) interface{} { return c.FS() }},

		// Private Infrastructure
		{"cacheStore", func(c *Container) interface{} { return c.cacheStore() }},
		{"clientProvider", func(c *Container) interface{} { return c.clientProvider() }},
		{"metadataCache", func(c *Container) interface{} { return c.metadataCache() }},
		{"metadataListingCache", func(c *Container) interface{} { return c.metadataListingCache() }},
		{"contentsCache", func(c *Container) interface{} { return c.contentsCache() }},
		{"pathIDCache", func(c *Container) interface{} { return c.pathIDCache() }},
		{"authCache", func(c *Container) interface{} { return c.authCache() }},
		{"accountCache", func(c *Container) interface{} { return c.accountCache() }},

		// Private Repositories
		{"metadataRepo", func(c *Container) interface{} { return c.metadata() }},
		{"contentsRepo", func(c *Container) interface{} { return c.contents() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(c)
			assert.NotNil(t, result, "Method %s() returned nil", tt.name)
		})
	}
}

func TestContainer_Laziness(t *testing.T) {
	setupTestEnv(t)
	c := NewContainer()

	tests := []struct {
		name string
		fn   func(*Container) interface{}
	}{
		{"Logger", func(c *Container) interface{} { return c.Logger() }},
		{"MetadataRepo", func(c *Container) interface{} { return c.metadata() }},
		{"EnvironmentService", func(c *Container) interface{} { return c.EnvironmentService() }},
		{"Config", func(c *Container) interface{} { return c.Config() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := tt.fn(c)
			s2 := tt.fn(c)
			assert.Same(t, s1, s2, "Method %s() should return a singleton instance", tt.name)
		})
	}
}
