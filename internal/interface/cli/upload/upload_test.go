package upload

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
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
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
func (m *MockContainer) FS() domainfs.Service    { return m.Called().Get(0).(domainfs.Service) }
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
func (m *MockContainer) State() domainstate.Service                          { return nil }
func (m *MockContainer) Drive() domaindrive.DriveService                     { return nil }
func (m *MockContainer) Account() domainaccount.Service                      { return nil }
func (m *MockContainer) Editor() domaineditor.Service                        { return nil }

type MockFSService struct {
	mock.Mock
}

func (m *MockFSService) Get(ctx context.Context, path string) (domainfs.Item, error) {
	return domainfs.Item{}, nil
}
func (m *MockFSService) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	return nil, nil
}
func (m *MockFSService) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	return nil, nil
}
func (m *MockFSService) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	return domainfs.Item{}, nil
}
func (m *MockFSService) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	return nil
}
func (m *MockFSService) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	return nil
}
func (m *MockFSService) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	return nil
}
func (m *MockFSService) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	return domainfs.Item{}, nil
}
func (m *MockFSService) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	args := m.Called(ctx, src, dst, opts)
	return args.Get(0).(domainfs.Item), args.Error(1)
}
func (m *MockFSService) Copy(ctx context.Context, src, dst string, opts domainfs.CopyOptions) error {
	return nil
}
func (m *MockFSService) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	return domainfs.Item{}, nil
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

func TestUploadCmd_Run(t *testing.T) {
	tests := []struct {
		name          string
		src           string
		dst           string
		mockError     error
		expectedError bool
	}{
		{
			name:          "success",
			src:           "local.txt",
			dst:           "/remote.txt",
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "error",
			src:           "local.txt",
			dst:           "/remote.txt",
			mockError:     errors.New("upload failed"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := new(MockFSService)
			mockLog := new(MockLogger)
			mockLogProvider := new(MockLogProvider)
			mockContainer := new(MockContainer)

			mockContainer.On("FS").Return(mockFS)
			mockContainer.On("Logger").Return(mockLogProvider)
			mockLogProvider.On("GetLogger", mock.Anything).Return(mockLog, nil)
			mockLogProvider.On("CreateLogger", mock.Anything).Return(mockLog, nil)
			mockLog.On("WithContext", mock.Anything).Return(mockLog)
			mockLog.On("With", mock.Anything).Return(mockLog)
			mockLog.On("Info", mock.Anything, mock.Anything).Return()
			mockLog.On("Error", mock.Anything, mock.Anything).Return()

			mockFS.On("Upload", mock.Anything, tt.src, tt.dst, mock.Anything).Return(domainfs.Item{}, tt.mockError)

			cmd := NewUploadCmd(mockContainer)
			opts := Options{
				Source:      tt.src,
				Destination: tt.dst,
				Stdout:      io.Discard,
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
