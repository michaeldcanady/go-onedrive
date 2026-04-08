package state

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/shared"
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

type mockLogger struct{}

func (m *mockLogger) Info(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Warn(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Error(msg string, kv ...logger.Field)          {}
func (m *mockLogger) Debug(msg string, kv ...logger.Field)          {}
func (m *mockLogger) SetLevel(level logger.Level)                   {}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger     { return m }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger { return m }

func TestBoltService_Operations(t *testing.T) {
	tmpDir := t.TempDir()
	mockEnv := new(MockEnvService)
	mockEnv.On("StateDir").Return(tmpDir, nil)
	l := &mockLogger{}

	service, err := NewBoltService(mockEnv, l)
	assert.NoError(t, err)
	defer service.Close()

	tests := []struct {
		name          string
		operation     func(s *BoltService) error
		check         func(t *testing.T, s *BoltService)
		expectedError error
	}{
		{
			name: "Initial state has default profile",
			check: func(t *testing.T, s *BoltService) {
				val, err := s.Get(KeyProfile)
				assert.NoError(t, err)
				assert.Equal(t, shared.DefaultProfileName, val)
			},
		},
		{
			name: "Set and Get Global",
			operation: func(s *BoltService) error {
				return s.Set(KeyDrive, "drive-123", ScopeGlobal)
			},
			check: func(t *testing.T, s *BoltService) {
				val, err := s.Get(KeyDrive)
				assert.NoError(t, err)
				assert.Equal(t, "drive-123", val)
			},
		},
		{
			name: "Set Session overrides Global",
			operation: func(s *BoltService) error {
				return s.Set(KeyDrive, "session-drive", ScopeSession)
			},
			check: func(t *testing.T, s *BoltService) {
				val, err := s.Get(KeyDrive)
				assert.NoError(t, err)
				assert.Equal(t, "session-drive", val)
			},
		},
		{
			name: "Clear removes from all scopes",
			operation: func(s *BoltService) error {
				return s.Clear(KeyDrive)
			},
			check: func(t *testing.T, s *BoltService) {
				_, err := s.Get(KeyDrive)
				assert.ErrorIs(t, err, ErrKeyNotFound)
			},
		},
		{
			name: "Get non-existent key errors",
			check: func(t *testing.T, s *BoltService) {
				_, err := s.Get(KeyAccessToken)
				assert.ErrorIs(t, err, ErrKeyNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.operation != nil {
				err := tt.operation(service)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				} else {
					assert.NoError(t, err)
				}
			}
			if tt.check != nil {
				tt.check(t, service)
			}
		})
	}
}

func TestBoltService_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	mockEnv := new(MockEnvService)
	mockEnv.On("StateDir").Return(tmpDir, nil)
	l := &mockLogger{}

	// Phase 1: Set data
	s1, _ := NewBoltService(mockEnv, l)
	_ = s1.Set(KeyAccessToken, "token-abc", ScopeGlobal)
	s1.Close()

	// Phase 2: Verify persistence
	tests := []struct {
		name     string
		key      Key
		expected string
	}{
		{
			name:     "Global token persisted",
			key:      KeyAccessToken,
			expected: "token-abc",
		},
	}

	s2, err := NewBoltService(mockEnv, l)
	assert.NoError(t, err)
	defer s2.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := s2.Get(tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestKey_String(t *testing.T) {
	tests := []struct {
		key      Key
		expected string
	}{
		{KeyProfile, "profile"},
		{KeyDrive, "drive"},
		{KeyAccessToken, "access_token"},
		{KeyConfigOverride, "config_override"},
		{Key(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.key.String())
		})
	}
}
