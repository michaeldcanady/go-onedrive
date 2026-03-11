package use

import (
	"context"
	"errors"
	"io"
	"testing"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal/config/domain"
	domainenv "github.com/michaeldcanady/go-onedrive/internal/core/env/domain"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal/editor/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal/profile/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockContainer struct {
	mock.Mock
}

func (m *MockContainer) Cache() pkgcache.Service { return nil }
func (m *MockContainer) FS() domainfs.Service    { return nil }
func (m *MockContainer) EnvironmentService() domainenv.EnvironmentService {
	return m.Called().Get(0).(domainenv.EnvironmentService)
}
func (m *MockContainer) Logger() domainlogger.LoggerService {
	return m.Called().Get(0).(domainlogger.LoggerService)
}
func (m *MockContainer) IgnoreMatcherFactory() domainfs.IgnoreMatcherFactory { return nil }
func (m *MockContainer) Auth() domainauth.AuthService                        { return nil }
func (m *MockContainer) Profile() domainprofile.ProfileService               { return nil }
func (m *MockContainer) Config() domainconfig.ConfigService                  { return nil }
func (m *MockContainer) State() domainstate.Service {
	return m.Called().Get(0).(domainstate.Service)
}
func (m *MockContainer) Drive() domaindrive.DriveService {
	return m.Called().Get(0).(domaindrive.DriveService)
}
func (m *MockContainer) Account() domainaccount.Service { return nil }
func (m *MockContainer) Editor() domaineditor.Service   { return nil }

type MockDriveService struct {
	mock.Mock
}

func (m *MockDriveService) ListDrives(ctx context.Context) ([]*domaindrive.Drive, error) {
	return nil, nil
}
func (m *MockDriveService) ResolvePersonalDrive(ctx context.Context) (*domaindrive.Drive, error) {
	return nil, nil
}
func (m *MockDriveService) ResolveDrive(ctx context.Context, target string) (*domaindrive.Drive, error) {
	args := m.Called(ctx, target)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domaindrive.Drive), args.Error(1)
}
func (m *MockDriveService) GetDrive(ctx context.Context, id string) (*domaindrive.Drive, error) {
	return nil, nil
}

type MockStateService struct {
	mock.Mock
}

func (m *MockStateService) Get(key domainstate.Key) (string, error) { return "", nil }
func (m *MockStateService) Set(key domainstate.Key, value string, scope domainstate.Scope) error {
	args := m.Called(key, value, scope)
	return args.Error(0)
}
func (m *MockStateService) Clear(key domainstate.Key) error { return nil }
func (m *MockStateService) GetDriveAlias(alias string) (string, error) {
	return "", nil
}
func (m *MockStateService) SetDriveAlias(alias, driveID string) error { return nil }
func (m *MockStateService) RemoveDriveAlias(alias string) error       { return nil }
func (m *MockStateService) ListDriveAliases() (map[string]string, error) {
	return nil, nil
}

type MockLogProvider struct {
	mock.Mock
}

func (m *MockLogProvider) CreateLogger(name string) (domainlogger.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(domainlogger.Logger), args.Error(1)
}
func (m *MockLogProvider) GetLogger(name string) (domainlogger.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(domainlogger.Logger), args.Error(1)
}
func (m *MockLogProvider) GetContextLogger(ctx context.Context, name string) (domainlogger.Logger, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(domainlogger.Logger), args.Error(1)
}
func (m *MockLogProvider) SetAllLevel(level string) { m.Called(level) }
func (m *MockLogProvider) SetLevel(id string, level string) error {
	args := m.Called(id, level)
	return args.Error(0)
}
func (m *MockLogProvider) RegisterProvider(t domainlogger.Type, factory domainlogger.LoggerProvider) {
	m.Called(t, factory)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...domainlogger.Field)  { m.Called(msg, fields) }
func (m *MockLogger) Error(msg string, fields ...domainlogger.Field) { m.Called(msg, fields) }
func (m *MockLogger) Debug(msg string, fields ...domainlogger.Field) { m.Called(msg, fields) }
func (m *MockLogger) Warn(msg string, fields ...domainlogger.Field)  { m.Called(msg, fields) }
func (m *MockLogger) SetLevel(level string)                          { m.Called(level) }
func (m *MockLogger) With(fields ...domainlogger.Field) domainlogger.Logger {
	args := m.Called(fields)
	if v := args.Get(0); v != nil {
		return v.(domainlogger.Logger)
	}
	return m
}
func (m *MockLogger) WithContext(ctx context.Context) domainlogger.Logger {
	args := m.Called(ctx)
	if v := args.Get(0); v != nil {
		return v.(domainlogger.Logger)
	}
	return m
}
func (m *MockLogger) GetContextLogger(ctx context.Context, name string) (domainlogger.Logger, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(domainlogger.Logger), args.Error(1)
}

// --- Tests ---

func TestDriveUseCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		target        string
		mockDrive     *domaindrive.Drive
		resolveError  error
		setError      error
		expectedError bool
	}{
		{
			name:   "success",
			target: "my-drive",
			mockDrive: &domaindrive.Drive{
				ID:   "123",
				Name: "My Drive",
			},
			expectedError: false,
		},
		{
			name:          "resolve error",
			target:        "bad",
			resolveError:  errors.New("not found"),
			expectedError: true,
		},
		{
			name:   "set state error",
			target: "my-drive",
			mockDrive: &domaindrive.Drive{
				ID:   "123",
				Name: "My Drive",
			},
			setError:      errors.New("save failed"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDriveSvc := new(MockDriveService)
			mockState := new(MockStateService)
			mockLog := new(MockLogger)
			mockLogProvider := new(MockLogProvider)
			mockContainer := new(MockContainer)

			mockContainer.On("Drive").Return(mockDriveSvc)
			mockContainer.On("State").Return(mockState)
			mockContainer.On("Logger").Return(mockLogProvider)
			mockLogProvider.On("GetLogger", mock.Anything).Return(mockLog, nil)
			mockLogProvider.On("CreateLogger", mock.Anything).Return(mockLog, nil)
			mockLog.On("WithContext", mock.Anything).Return(mockLog)
			mockLog.On("With", mock.Anything).Return(mockLog)
			mockLog.On("Info", mock.Anything, mock.Anything).Return()
			mockLog.On("Error", mock.Anything, mock.Anything).Return()

			mockDriveSvc.On("ResolveDrive", mock.Anything, tt.target).Return(tt.mockDrive, tt.resolveError)
			if tt.resolveError == nil && tt.mockDrive != nil {
				mockState.On("Set", domainstate.KeyDrive, tt.mockDrive.ID, domainstate.ScopeGlobal).Return(tt.setError)
			}

			cmd := NewUseCmd(mockContainer)
			opts := Options{
				DriveIDOrAlias: tt.target,
				Stdout:         io.Discard,
			}

			err := cmd.Run(context.Background(), opts)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
