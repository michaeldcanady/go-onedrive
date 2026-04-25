package config

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProfileProvider struct {
	mock.Mock
}

func (m *mockProfileProvider) GetActive(ctx context.Context) (profile.Profile, error) {
	args := m.Called(ctx)
	return args.Get(0).(profile.Profile), args.Error(1)
}

type mockLoggerSvc struct {
	mock.Mock
}

func (m *mockLoggerSvc) Info(msg string, kv ...logger.Field)           {}
func (m *mockLoggerSvc) Warn(msg string, kv ...logger.Field)           {}
func (m *mockLoggerSvc) Error(msg string, kv ...logger.Field)          {}
func (m *mockLoggerSvc) Debug(msg string, kv ...logger.Field)          {}
func (m *mockLoggerSvc) SetLevel(level logger.Level)                   {}
func (m *mockLoggerSvc) With(fields ...logger.Field) logger.Logger     { return m }
func (m *mockLoggerSvc) WithContext(ctx context.Context) logger.Logger { return m }

func TestConfigService_GetPath(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(m *mockProfileProvider)
		override string
		want     string
		wantOk   bool
	}{
		{
			name: "override takes precedence",
			setup: func(m *mockProfileProvider) {
				m.On("GetActive", mock.Anything).Return(profile.Profile{ConfigPath: "/profile/path"}, nil)
			},
			override: "/override/path",
			want:     "/override/path",
			wantOk:   true,
		},
		{
			name: "profile path used if no override",
			setup: func(m *mockProfileProvider) {
				m.On("GetActive", mock.Anything).Return(profile.Profile{ConfigPath: "/profile/path"}, nil)
			},
			want:   "/profile/path",
			wantOk: true,
		},
		{
			name: "no path if both missing",
			setup: func(m *mockProfileProvider) {
				m.On("GetActive", mock.Anything).Return(profile.Profile{}, errors.New("no active profile"))
			},
			want:   "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProfile := new(mockProfileProvider)
			tt.setup(mProfile)

			svc := NewConfigService(mProfile, &mockLoggerSvc{})
			if tt.override != "" {
				_ = svc.SetOverride(context.Background(), tt.override)
			}

			path, ok := svc.GetPath(context.Background())
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.want, path)
		})
	}
}

func TestConfigService_UpdateConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create initial config file
	initialCfg := []byte("auth:\n  provider: microsoft\n")
	err := os.WriteFile(configPath, initialCfg, 0644)
	assert.NoError(t, err)

	svc := NewConfigService(nil, &mockLoggerSvc{})
	_ = svc.SetOverride(context.Background(), configPath)

	t.Run("update success", func(t *testing.T) {
		err := svc.UpdateConfig(context.Background(), "auth.provider", "google")
		assert.NoError(t, err)

		cfg, err := svc.GetConfig(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "google", cfg.Auth.Provider)
	})

	t.Run("update failure - invalid key", func(t *testing.T) {
		err := svc.UpdateConfig(context.Background(), "invalid.key", "val")
		assert.Error(t, err)
	})
}
