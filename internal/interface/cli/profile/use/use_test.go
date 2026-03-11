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
func (m *MockContainer) Profile() domainprofile.ProfileService {
	return m.Called().Get(0).(domainprofile.ProfileService)
}
func (m *MockContainer) Config() domainconfig.ConfigService { return nil }
func (m *MockContainer) State() domainstate.Service {
	return m.Called().Get(0).(domainstate.Service)
}
func (m *MockContainer) Drive() domaindrive.DriveService { return nil }
func (m *MockContainer) Account() domainaccount.Service  { return nil }
func (m *MockContainer) Editor() domaineditor.Service    { return nil }

type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) Get(ctx context.Context, name string) (domainprofile.Profile, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return domainprofile.Profile{}, args.Error(1)
	}
	return args.Get(0).(domainprofile.Profile), args.Error(1)
}
func (m *MockProfileService) List(ctx context.Context) ([]domainprofile.Profile, error) {
	return nil, nil
}
func (m *MockProfileService) Create(ctx context.Context, name string) (domainprofile.Profile, error) {
	return domainprofile.Profile{}, nil
}
func (m *MockProfileService) Delete(ctx context.Context, name string) error { return nil }
func (m *MockProfileService) Exists(ctx context.Context, name string) (bool, error) {
	return false, nil
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

func TestProfileUseCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		profileName   string
		mockProfile   domainprofile.Profile
		getError      error
		setError      error
		expectedError bool
	}{
		{
			name:          "success",
			profileName:   "personal",
			mockProfile:   domainprofile.Profile{Name: "personal"},
			expectedError: false,
		},
		{
			name:          "profile not found",
			profileName:   "missing",
			getError:      errors.New("not found"),
			expectedError: true,
		},
		{
			name:          "set state error",
			profileName:   "personal",
			mockProfile:   domainprofile.Profile{Name: "personal"},
			setError:      errors.New("save failed"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProfileSvc := new(MockProfileService)
			mockState := new(MockStateService)
			mockLog := new(MockLogger)
			mockLogProvider := new(MockLogProvider)
			mockContainer := new(MockContainer)

			mockContainer.On("Profile").Return(mockProfileSvc)
			mockContainer.On("State").Return(mockState)
			mockContainer.On("Logger").Return(mockLogProvider)
			mockLogProvider.On("GetLogger", mock.Anything).Return(mockLog, nil)
			mockLogProvider.On("CreateLogger", mock.Anything).Return(mockLog, nil)
			mockLog.On("WithContext", mock.Anything).Return(mockLog)
			mockLog.On("With", mock.Anything).Return(mockLog)
			mockLog.On("Info", mock.Anything, mock.Anything).Return()
			mockLog.On("Error", mock.Anything, mock.Anything).Return()

			if tt.getError != nil {
				mockProfileSvc.On("Get", mock.Anything, tt.profileName).Return(domainprofile.Profile{}, tt.getError)
			} else {
				mockProfileSvc.On("Get", mock.Anything, tt.profileName).Return(tt.mockProfile, nil)
			}

			if tt.getError == nil {
				mockState.On("Set", domainstate.KeyProfile, tt.mockProfile.Name, domainstate.ScopeGlobal).Return(tt.setError)
			}

			cmd := NewUseCmd(mockContainer)
			opts := Options{
				Name:   tt.profileName,
				Stdout: io.Discard,
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
