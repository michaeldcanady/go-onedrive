package edit

import (
	"bytes"
	"context"
	"io"
	"strings"
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
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
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
func (m *MockContainer) File() file.FileService                { return nil }
func (m *MockContainer) State() state.Service                  { return nil }
func (m *MockContainer) Drive() drive.DriveService             { return nil }
func (m *MockContainer) Account() account.Service              { return nil }
func (m *MockContainer) Editor() domaineditor.Service {
	args := m.Called()
	return args.Get(0).(domaineditor.Service)
}

type MockFSService struct {
	mock.Mock
}

func (m *MockFSService) Get(ctx context.Context, path string) (domainfs.Item, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(domainfs.Item), args.Error(1)
}
func (m *MockFSService) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	return nil, nil
}
func (m *MockFSService) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	return domainfs.Item{}, nil
}
func (m *MockFSService) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, path, opts)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
func (m *MockFSService) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	args := m.Called(ctx, path, r, opts)
	return args.Get(0).(domainfs.Item), args.Error(1)
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
func (m *MockFSService) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	return domainfs.Item{}, nil
}

type MockEnvironmentService struct {
	mock.Mock
}

func (m *MockEnvironmentService) CacheDir() (string, error)   { return "", nil }
func (m *MockEnvironmentService) ConfigDir() (string, error)  { return "", nil }
func (m *MockEnvironmentService) DataDir() (string, error)    { return "", nil }
func (m *MockEnvironmentService) EnsureAll() error            { return nil }
func (m *MockEnvironmentService) InstallDir() (string, error) { return "", nil }
func (m *MockEnvironmentService) IsLinux() bool               { return m.Called().Bool(0) }
func (m *MockEnvironmentService) IsMac() bool                 { return m.Called().Bool(0) }
func (m *MockEnvironmentService) IsWindows() bool             { return m.Called().Bool(0) }
func (m *MockEnvironmentService) LogDir() (string, error)     { return "", nil }
func (m *MockEnvironmentService) Name() string                { return "odc" }
func (m *MockEnvironmentService) OS() string                  { return "linux" }
func (m *MockEnvironmentService) TempDir() (string, error)    { return "", nil }
func (m *MockEnvironmentService) StateDir() (string, error)   { return "", nil }
func (m *MockEnvironmentService) OutputDestination() (infralogging.OutputDestination, error) {
	return 0, nil
}
func (m *MockEnvironmentService) LogLevel() (string, error) { return "info", nil }
func (m *MockEnvironmentService) Shell() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *MockEnvironmentService) Editor() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *MockEnvironmentService) Visual() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockLogProvider struct {
	mock.Mock
}

func (m *MockLogProvider) CreateLogger(name string) (infralogging.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(infralogging.Logger), args.Error(1)
}
func (m *MockLogProvider) GetLogger(name string) (infralogging.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(infralogging.Logger), args.Error(1)
}
func (m *MockLogProvider) SetAllLevel(level string) {}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...infralogging.Field)  { m.Called(msg, fields) }
func (m *MockLogger) Error(msg string, fields ...infralogging.Field) { m.Called(msg, fields) }
func (m *MockLogger) Debug(msg string, fields ...infralogging.Field) { m.Called(msg, fields) }
func (m *MockLogger) Warn(msg string, fields ...infralogging.Field)  { m.Called(msg, fields) }
func (m *MockLogger) SetLevel(level string)                          {}
func (m *MockLogger) With(fields ...infralogging.Field) infralogging.Logger {
	return m
}
func (m *MockLogger) WithContext(ctx context.Context) infralogging.Logger {
	return m
}

type MockEditor struct {
	mock.Mock
}

func (m *MockEditor) Launch(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockEditor) LaunchTempFile(prefix, suffix string, reader io.Reader) ([]byte, string, error) {
	_, _ = io.ReadAll(reader) // Consume reader for hash calculation
	args := m.Called(prefix, suffix, reader)
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}

func (m *MockEditor) WithIO(stdin io.Reader, stdout, stderr io.Writer) domaineditor.Service {
	m.Called(stdin, stdout, stderr)
	return m
}

// --- Tests ---

func TestName(t *testing.T) {
	assert.Equal(t, "file", Name("/path/to/file.txt"))
	assert.Equal(t, "data", Name("data.json"))
	assert.Equal(t, "noext", Name("noext"))
}

func TestCreateEditCmd_Flags(t *testing.T) {
	mockContainer := new(MockContainer)
	cmd := CreateEditCmd(mockContainer)

	assert.Equal(t, "edit [path]", cmd.Use)
}

func TestEditCmd_Run_NoChanges(t *testing.T) {
	mockContainer := new(MockContainer)
	mockFS := new(MockFSService)
	mockEnv := new(MockEnvironmentService)
	mockLogProvider := new(MockLogProvider)
	mockLogger := new(MockLogger)
	mockEditor := new(MockEditor)

	mockContainer.On("FS").Return(mockFS)
	mockContainer.On("EnvironmentService").Return(mockEnv)
	mockContainer.On("Logger").Return(mockLogProvider)
	mockLogProvider.On("GetLogger", mock.Anything).Return(mockLogger, nil)
	mockLogProvider.On("CreateLogger", mock.Anything).Return(mockLogger, nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	// Mock ReadFile
	content := "original content"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor
	mockEditor.On("WithIO", mock.Anything, mock.Anything, mock.Anything).Return(mockEditor)
	mockEditor.On("LaunchTempFile", mock.Anything, ".txt", mock.Anything).
		Return([]byte(content), "temp-path", nil)

	buf := new(bytes.Buffer)
	opts := Options{Path: "/file.txt", Stdout: buf}

	editCmd := NewEditCmd(mockContainer).WithLogger(mockLogger).WithEditor(mockEditor)
	err := editCmd.Run(context.Background(), opts)

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No changes detected.")
}

func TestEditCmd_Run_WithChanges(t *testing.T) {
	mockContainer := new(MockContainer)
	mockFS := new(MockFSService)
	mockEnv := new(MockEnvironmentService)
	mockLogProvider := new(MockLogProvider)
	mockLogger := new(MockLogger)
	mockEditor := new(MockEditor)

	mockContainer.On("FS").Return(mockFS)
	mockContainer.On("EnvironmentService").Return(mockEnv)
	mockContainer.On("Logger").Return(mockLogProvider)
	mockLogProvider.On("GetLogger", mock.Anything).Return(mockLogger, nil)
	mockLogProvider.On("CreateLogger", mock.Anything).Return(mockLogger, nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	// Mock ReadFile
	content := "original content"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor
	mockEditor.On("WithIO", mock.Anything, mock.Anything, mock.Anything).Return(mockEditor)
	mockEditor.On("LaunchTempFile", mock.Anything, ".txt", mock.Anything).
		Return([]byte("new content"), "temp-path", nil)

	// Mock WriteFile: should be called because content changed
	mockFS.On("WriteFile", mock.Anything, "/file.txt", mock.Anything, mock.Anything).
		Return(domainfs.Item{Name: "file.txt"}, nil)

	buf := new(bytes.Buffer)
	opts := Options{Path: "/file.txt", Stdout: buf}

	editCmd := NewEditCmd(mockContainer).WithLogger(mockLogger).WithEditor(mockEditor)
	err := editCmd.Run(context.Background(), opts)

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "updated successfully")
	mockFS.AssertExpectations(t)
}

func TestEditCmd_Run_WithMockEditor(t *testing.T) {
	mockContainer := new(MockContainer)
	mockFS := new(MockFSService)
	mockLogger := new(MockLogger)
	mockEditor := new(MockEditor)

	mockContainer.On("FS").Return(mockFS)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	// Mock ReadFile
	content := "original"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor: returns modified content
	mockEditor.On("WithIO", mock.Anything, mock.Anything, mock.Anything).Return(mockEditor)
	mockEditor.On("LaunchTempFile", mock.Anything, ".txt", mock.Anything).
		Return([]byte("modified"), "temp-path", nil)

	// Mock WriteFile: should be called because content changed
	mockFS.On("WriteFile", mock.Anything, "/file.txt", mock.Anything, mock.Anything).
		Return(domainfs.Item{Name: "file.txt"}, nil)

	buf := new(bytes.Buffer)
	opts := Options{Path: "/file.txt", Stdout: buf}

	editCmd := NewEditCmd(mockContainer).WithLogger(mockLogger).WithEditor(mockEditor)
	err := editCmd.Run(context.Background(), opts)

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "updated successfully")
	mockEditor.AssertExpectations(t)
	mockFS.AssertExpectations(t)
}
