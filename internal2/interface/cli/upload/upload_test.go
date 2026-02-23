package upload

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenvironment "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type MockContainer struct {
	mock.Mock
}

func (m *MockContainer) Cache() domaincache.Service2               { return nil }
func (m *MockContainer) FS() domainfs.Service                      { return m.Called().Get(0).(domainfs.Service) }
func (m *MockContainer) EnvironmentService() domainenvironment.EnvironmentService {
	return nil
}
func (m *MockContainer) Logger() domainlogger.LoggerService {
	return m.Called().Get(0).(domainlogger.LoggerService)
}
func (m *MockContainer) Auth() domainauth.AuthService       { return nil }
func (m *MockContainer) Profile() domainprofile.ProfileService { return nil }
func (m *MockContainer) Config() config.ConfigService       { return nil }
func (m *MockContainer) File() file.FileService             { return nil }
func (m *MockContainer) State() state.Service               { return nil }
func (m *MockContainer) Drive() drive.DriveService           { return nil }
func (m *MockContainer) Account() account.Service           { return nil }
func (m *MockContainer) Editor() domaineditor.Service        { return nil }

type MockFSService struct {
	mock.Mock
}

func (m *MockFSService) Get(ctx context.Context, path string) (domainfs.Item, error) {
	return domainfs.Item{}, nil
}
func (m *MockFSService) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	return nil, nil
}
func (m *MockFSService) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	return domainfs.Item{}, nil
}
func (m *MockFSService) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	return nil, nil
}
func (m *MockFSService) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	args := m.Called(ctx, path, r, opts)
	return args.Get(0).(domainfs.Item), args.Error(1)
}
func (m *MockFSService) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	return nil
}
func (m *MockFSService) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	return nil
}
func (m *MockFSService) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	return nil
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
	m.Called(fields)
	return m
}
func (m *MockLogger) WithContext(ctx context.Context) logging.Logger {
	m.Called(ctx)
	return m
}

func TestUploadCmd_Run(t *testing.T) {
	// We need a real file to open, or we can mock os.OpenFile (which is harder).
	// Let's create a temporary file for the test.
	tmpFile, err := os.CreateTemp("", "testupload-*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("hello world")
	tmpFile.Close()

	t.Run("Success", func(t *testing.T) {
		mockContainer := new(MockContainer)
		mockFS := new(MockFSService)
		mockLogProvider := new(MockLogProvider)
		mockLogger := new(MockLogger)

		mockContainer.On("FS").Return(mockFS)
		mockContainer.On("Logger").Return(mockLogProvider)
		mockLogProvider.On("GetLogger", mock.Anything).Return(mockLogger, nil)
		mockLogProvider.On("CreateLogger", mock.Anything).Return(mockLogger, nil)

		mockLogger.On("Info", mock.Anything, mock.Anything).Return()
		mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

		mockFS.On("WriteFile", mock.Anything, "/remote.txt", mock.Anything, domainfs.WriteOptions{Overwrite: true}).
			Return(domainfs.Item{Name: "remote.txt"}, nil)

		cmd := NewUploadCmd(mockContainer).WithLogger(mockLogger)
		opts := Options{
			Source:      tmpFile.Name(),
			Destination: "/remote.txt",
			Overwrite:   true,
			Stdout:      io.Discard,
		}

		err := cmd.Run(context.Background(), opts)
		assert.NoError(t, err)
		mockFS.AssertExpectations(t)
	})
}
func TestUploadCmd_ResolveDestination(t *testing.T) {
	cmd := &UploadCmd{}

	tests := []struct {
		src      string
		dst      string
		expected string
	}{
		{"local.txt", "/remote.txt", "/remote.txt"},
		{"local.txt", "/folder/", "/folder/local.txt"},
		{"/path/to/local.txt", "/folder/", "/folder/local.txt"},
	}

	for _, tt := range tests {
		result := cmd.resolveDestination(tt.src, tt.dst)
		assert.Equal(t, tt.expected, result)
	}
}

func TestCreateUploadCmd_Flags(t *testing.T) {
	mockContainer := new(MockContainer)
	cmd := CreateUploadCmd(mockContainer)

	assert.Equal(t, "upload [src] [dst]", cmd.Use)
	assert.NotNil(t, cmd.Flags().Lookup("force"))
}
