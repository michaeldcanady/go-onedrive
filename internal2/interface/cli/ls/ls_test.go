package ls

import (
	"context"
	"io"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenvironment "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockContainer struct {
	mock.Mock
}

func (m *MockContainer) Cache() domaincache.Service2            { return nil }
func (m *MockContainer) FS() domainfs.Service                   { return m.Called().Get(0).(domainfs.Service) }
func (m *MockContainer) EnvironmentService() domainenvironment.EnvironmentService {
	return m.Called().Get(0).(domainenvironment.EnvironmentService)
}
func (m *MockContainer) Logger() domainlogger.LoggerService {
	return m.Called().Get(0).(domainlogger.LoggerService)
}
func (m *MockContainer) Auth() domainauth.AuthService          { return nil }
func (m *MockContainer) Profile() domainprofile.ProfileService { return nil }
func (m *MockContainer) Config() config.ConfigService          { return nil }
func (m *MockContainer) State() state.Service                  { return m.Called().Get(0).(state.Service) }
func (m *MockContainer) Drive() drive.DriveService             { return nil }
func (m *MockContainer) Account() account.Service              { return nil }
func (m *MockContainer) Editor() domaineditor.Service          { return nil }

type MockFSService struct {
	mock.Mock
}

func (m *MockFSService) Get(ctx context.Context, path string) (domainfs.Item, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(domainfs.Item), args.Error(1)
}

func (m *MockFSService) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	args := m.Called(ctx, path, opts)
	return args.Get(0).([]domainfs.Item), args.Error(1)
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
	return domainfs.Item{}, nil
}

type MockLogProvider struct {
	mock.Mock
}

func (m *MockLogProvider) CreateLogger(name string) (logging.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(logging.Logger), args.Error(1)
}

func (m *MockLogProvider) GetLogger(name string) (logging.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(logging.Logger), args.Error(1)
}

func (m *MockLogProvider) SetAllLevel(level string) {
	m.Called(level)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...logging.Field)  { m.Called(msg, fields) }
func (m *MockLogger) Error(msg string, fields ...logging.Field) { m.Called(msg, fields) }
func (m *MockLogger) Debug(msg string, fields ...logging.Field) { m.Called(msg, fields) }
func (m *MockLogger) Warn(msg string, fields ...logging.Field)  { m.Called(msg, fields) }
func (m *MockLogger) SetLevel(level string)                     { m.Called(level) }
func (m *MockLogger) With(fields ...logging.Field) logging.Logger {
	args := m.Called(fields)
	if v := args.Get(0); v != nil {
		return v.(logging.Logger)
	}
	return m
}
func (m *MockLogger) WithContext(ctx context.Context) logging.Logger {
	args := m.Called(ctx)
	if v := args.Get(0); v != nil {
		return v.(logging.Logger)
	}
	return m
}

// --- Tests ---

func TestCreateLSCmd_Flags(t *testing.T) {
	mockContainer := new(MockContainer)
	cmd := CreateLSCmd(mockContainer)

	assert.Equal(t, "ls [path]", cmd.Use)

	// Check flags
	f := cmd.Flags()
	assert.NotNil(t, f.Lookup("all"))
	assert.NotNil(t, f.Lookup("format"))
	assert.NotNil(t, f.Lookup("sort"))
	assert.NotNil(t, f.Lookup("folders-only"))
	assert.NotNil(t, f.Lookup("files-only"))
	assert.NotNil(t, f.Lookup("recursive"))
	assert.NotNil(t, f.Lookup("long"))
	assert.NotNil(t, f.Lookup("tree"))
}
