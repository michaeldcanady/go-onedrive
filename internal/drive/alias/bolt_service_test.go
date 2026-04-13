package alias

import (
	"context"
	"testing"

	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
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
	mockEnv.On("StateDir").Return(tmpDir, nil)
	mockLog := new(MockLogger)
	mockLog.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLog.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLog.On("Warn", mock.Anything, mock.Anything).Maybe()
	mockLog.On("Error", mock.Anything, mock.Anything).Maybe()

	service, err := NewBoltService(mockEnv, mockLog)
	assert.NoError(t, err)
	defer service.Close()

	tests := []struct {
		name      string
		operation func(t *testing.T, s *BoltService)
	}{
		{
			name: "Set and Get alias",
			operation: func(t *testing.T, s *BoltService) {
				driveID := "drive-123"
				aliasName := "my-drive"
				err := s.SetAlias(driveID, aliasName)
				assert.NoError(t, err)

				gotDriveID, err := s.GetDriveIDByAlias(aliasName)
				assert.NoError(t, err)
				assert.Equal(t, driveID, gotDriveID)

				gotAlias, err := s.GetAliasByDriveID(driveID)
				assert.NoError(t, err)
				assert.Equal(t, aliasName, gotAlias)
			},
		},
		{
			name: "Get non-existent alias returns ErrAliasNotFound",
			operation: func(t *testing.T, s *BoltService) {
				_, err := s.GetDriveIDByAlias("non-existent")
				assert.ErrorIs(t, err, coreerrors.CodeNotFound)
			},
		},
		{
			name: "List aliases",
			operation: func(t *testing.T, s *BoltService) {
				_ = s.SetAlias("d1", "a1")
				_ = s.SetAlias("d2", "a2")

				aliases, err := s.ListAliases()
				assert.NoError(t, err)
				assert.Contains(t, aliases, "a1")
				assert.Contains(t, aliases, "a2")
				assert.Equal(t, "d1", aliases["a1"])
				assert.Equal(t, "d2", aliases["a2"])
			},
		},
		{
			name: "Delete alias",
			operation: func(t *testing.T, s *BoltService) {
				_ = s.SetAlias("d3", "a3")
				err := s.DeleteAlias("a3")
				assert.NoError(t, err)

				_, err = s.GetDriveIDByAlias("a3")
				assert.ErrorIs(t, err, coreerrors.CodeNotFound)
			},
		},
		{
			name: "Delete non-existent alias does not error",
			operation: func(t *testing.T, s *BoltService) {
				err := s.DeleteAlias("missing")
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.operation(t, service)
		})
	}
}
