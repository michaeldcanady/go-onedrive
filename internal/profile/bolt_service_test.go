package profile

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/shared"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEnvService is a mock implementation of environment.Service.
type MockEnvService struct {
	mock.Mock
}

func (m *MockEnvService) StateDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) CacheDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) ConfigDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) DataDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) EnsureAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockEnvService) InstallDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) IsLinux() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockEnvService) IsMac() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockEnvService) IsWindows() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockEnvService) LogDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEnvService) OS() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEnvService) TempDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) Shell() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) Editor() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) Visual() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockEnvService) LogLevel() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEnvService) LogOutput() string {
	args := m.Called()
	return args.String(0)
}

// MockStateService is a mock implementation of state.Service.
type MockStateService struct {
	mock.Mock
}

func (m *MockStateService) Get(key state.Key) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockStateService) Set(key state.Key, value string, scope state.Scope) error {
	args := m.Called(key, value, scope)
	return args.Error(0)
}

func (m *MockStateService) Clear(key state.Key) error {
	args := m.Called(key)
	return args.Error(0)
}

// MockLogger is a mock implementation of logger.Logger.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) Warn(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) Error(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) Debug(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) SetLevel(level logger.Level) {
	m.Called(level)
}

func (m *MockLogger) With(fields ...logger.Field) logger.Logger {
	return m
}

func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	return m
}

func TestBoltService_LifeCycle(t *testing.T) {
	tmpDir := t.TempDir()
	mockEnv := new(MockEnvService)
	mockEnv.On("ConfigDir").Return(tmpDir, nil)
	mockState := new(MockStateService)
	mockLog := new(MockLogger)
	mockLog.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLog.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLog.On("Warn", mock.Anything, mock.Anything).Maybe()
	mockLog.On("Error", mock.Anything, mock.Anything).Maybe()

	service, err := NewBoltService(mockEnv, mockState, mockLog)
	assert.NoError(t, err)
	defer service.Close()

	ctx := context.Background()

	tests := []struct {
		name      string
		operation func(t *testing.T, s *BoltService)
	}{
		{
			name: "Default profile exists initially",
			operation: func(t *testing.T, s *BoltService) {
				exists, err := s.Exists(ctx, shared.DefaultProfileName)
				assert.NoError(t, err)
				assert.True(t, exists)
			},
		},
		{
			name: "Create new profile",
			operation: func(t *testing.T, s *BoltService) {
				name := "test-profile"
				p, err := s.Create(ctx, name)
				assert.NoError(t, err)
				assert.Equal(t, name, p.Name)

				exists, _ := s.Exists(ctx, name)
				assert.True(t, exists)
			},
		},
		{
			name: "Create existing profile errors",
			operation: func(t *testing.T, s *BoltService) {
				_, err := s.Create(ctx, "test-profile")
				assert.ErrorIs(t, err, ErrProfileAlreadyExists)
			},
		},
		{
			name: "Get profile",
			operation: func(t *testing.T, s *BoltService) {
				p, err := s.Get(ctx, "test-profile")
				assert.NoError(t, err)
				assert.Equal(t, "test-profile", p.Name)
			},
		},
		{
			name: "Get non-existent profile errors",
			operation: func(t *testing.T, s *BoltService) {
				_, err := s.Get(ctx, "missing")
				assert.ErrorIs(t, err, ErrProfileNotFound)
			},
		},
		{
			name: "Update profile",
			operation: func(t *testing.T, s *BoltService) {
				p, _ := s.Get(ctx, "test-profile")
				p.ConfigPath = "/updated/path.yaml"
				err := s.Update(ctx, p)
				assert.NoError(t, err)

				got, _ := s.Get(ctx, "test-profile")
				assert.Equal(t, "/updated/path.yaml", got.ConfigPath)
			},
		},
		{
			name: "List profiles",
			operation: func(t *testing.T, s *BoltService) {
				profiles, err := s.List(ctx)
				assert.NoError(t, err)
				assert.Len(t, profiles, 2)
			},
		},
		{
			name: "SetActive / GetActive integration",
			operation: func(t *testing.T, s *BoltService) {
				name := "active-one"
				_, _ = s.Create(ctx, name)

				mockState.On("Set", state.KeyProfile, name, state.ScopeGlobal).Return(nil).Once()
				err := s.SetActive(ctx, name, state.ScopeGlobal)
				assert.NoError(t, err)

				mockState.On("Get", state.KeyProfile).Return(name, nil).Once()
				active, err := s.GetActive(ctx)
				assert.NoError(t, err)
				assert.Equal(t, name, active.Name)
			},
		},
		{
			name: "Delete profile",
			operation: func(t *testing.T, s *BoltService) {
				err := s.Delete(ctx, "test-profile")
				assert.NoError(t, err)

				exists, _ := s.Exists(ctx, "test-profile")
				assert.False(t, exists)
			},
		},
		{
			name: "Delete default profile errors",
			operation: func(t *testing.T, s *BoltService) {
				err := s.Delete(ctx, shared.DefaultProfileName)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.operation(t, service)
		})
	}
	mockState.AssertExpectations(t)
}
