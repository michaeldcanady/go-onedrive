package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestYAMLService_SaveAndGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "odc-config-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	l := &mockLogger{}
	svc := NewYAMLService(nil, l)
	ctx := context.Background()

	profile := "test-profile"
	configPath := filepath.Join(tmpDir, "config.yaml")

	// 1. Register path
	err = svc.AddPath(profile, configPath)
	assert.NoError(t, err)

	p, ok := svc.GetPath(profile)
	assert.True(t, ok)
	assert.Equal(t, configPath, p)

	// 2. Save config
	cfg := Config{
		Auth: AuthenticationConfig{
			Provider:    "microsoft",
			ClientID:    "new-client-id",
			RedirectURI: "http://localhost:1234",
		},
	}

	err = svc.SaveConfig(ctx, profile, cfg)
	assert.NoError(t, err)

	// 3. Verify file exists and content
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// 4. Get config back
	loadedCfg, err := svc.GetConfig(ctx, profile)
	assert.NoError(t, err)
	assert.Equal(t, "new-client-id", loadedCfg.Auth.ClientID)
	assert.Equal(t, "http://localhost:1234", loadedCfg.Auth.RedirectURI)
	// Verify defaults were merged for missing fields
	assert.Equal(t, "common", loadedCfg.Auth.TenantID)
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, kv ...logger.Field)                {}
func (m *mockLogger) Warn(msg string, kv ...logger.Field)                {}
func (m *mockLogger) Error(msg string, kv ...logger.Field)               {}
func (m *mockLogger) Debug(msg string, kv ...logger.Field)               {}
func (m *mockLogger) SetLevel(level string)                              {}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger          { return m }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger       { return m }
