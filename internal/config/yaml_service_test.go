package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestYAMLService_SaveAndGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "odc-config-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	l := &mockLogger{}
	st := &mockStateService{data: make(map[string]string)}
	svc := NewYAMLService(nil, st, l)
	ctx := context.Background()

	configPath := filepath.Join(tmpDir, "config.yaml")

	// 1. Register path override via state
	st.data["config_override"] = configPath

	p, ok := svc.GetPath(ctx)
	assert.True(t, ok)
	assert.Equal(t, configPath, p)

	// 2. Save config
	cfg := Config{
		Auth: AuthenticationConfig{
			Provider:    AuthProviderMicrosoft,
			ClientID:    "new-client-id",
			RedirectURI: "http://localhost:1234",
		},
	}

	err = svc.SaveConfig(ctx, cfg)
	assert.NoError(t, err)

	// 3. Verify file exists and content
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// 4. Get config back
	loadedCfg, err := svc.GetConfig(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "new-client-id", loadedCfg.Auth.ClientID)
	assert.Equal(t, "http://localhost:1234", loadedCfg.Auth.RedirectURI)
	// Verify defaults were merged for missing fields
	assert.Equal(t, "common", loadedCfg.Auth.TenantID)
}

func TestYAMLService_InvalidYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "odc-config-invalid-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err = os.WriteFile(configPath, []byte("auth: { invalid: yaml"), 0o644)
	assert.NoError(t, err)

	l := &mockLogger{}
	st := &mockStateService{data: make(map[string]string)}
	st.data["config_override"] = configPath
	svc := NewYAMLService(nil, st, l)
	ctx := context.Background()

	_, err = svc.GetConfig(ctx)
	assert.Error(t, err)

	var appErr *errors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, errors.CodeInvalidConfig, appErr.Code)
	assert.Contains(t, appErr.SafeMsg, "format is invalid")
	assert.Contains(t, appErr.Hint, "Ensure the YAML format is correct")
}

type mockStateService struct {
	data map[string]string
}

func (m *mockStateService) Get(key state.Key) (string, error) {
	return m.data["config_override"], nil
}

func (m *mockStateService) Set(key state.Key, value string, scope state.Scope) error {
	m.data["config_override"] = value
	return nil
}

func (m *mockStateService) Clear(key state.Key) error { return nil }

func (m *mockStateService) GetDriveAlias(alias string) (string, error) { return "", nil }
func (m *mockStateService) SetDriveAlias(alias, driveID string) error  { return nil }
func (m *mockStateService) RemoveDriveAlias(alias string) error        { return nil }
func (m *mockStateService) ListDriveAliases() (map[string]string, error) {
	return make(map[string]string), nil
}

func (m *mockStateService) GetDriveAliasByDriveID(driveID string) (string, error) {
	return "", nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Warn(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Error(msg string, kv ...logger.Field)          {}
func (m *mockLogger) Debug(msg string, kv ...logger.Field)          {}
func (m *mockLogger) SetLevel(level logger.Level)                   {}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger     { return m }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger { return m }
