package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

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
	svc := NewConfigService(nil, st, l)
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
			Provider:    "microsoft",
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

type mockStateService struct {
	data map[string]string
}

func (m *mockStateService) Get(key state.Key) (string, error) {
	return m.data[key.String()], nil
}

func (m *mockStateService) Set(key state.Key, value string, scope state.Scope) error {
	m.data[key.String()] = value
	return nil
}

func (m *mockStateService) Clear(key state.Key) error {
	delete(m.data, key.String())
	return nil
}

func (m *mockStateService) GetScoped(bucket, key string) (string, error) {
	return m.data[bucket+"/"+key], nil
}

func (m *mockStateService) SetScoped(bucket, key, value string, scope state.Scope) error {
	m.data[bucket+"/"+key] = value
	return nil
}

func (m *mockStateService) ClearScoped(bucket, key string) error {
	delete(m.data, bucket+"/"+key)
	return nil
}

func (m *mockStateService) ListScoped(bucket string) ([]string, error) {
	var keys []string
	prefix := bucket + "/"
	for k := range m.data {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k[len(prefix):])
		}
	}
	return keys, nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Warn(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Error(msg string, kv ...logger.Field)          {}
func (m *mockLogger) Debug(msg string, kv ...logger.Field)          {}
func (m *mockLogger) SetLevel(level logger.Level)                   {}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger     { return m }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger { return m }
