package edit

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	domaincache "github.com/michaeldcanady/go-onedrive/internal/cache/domain"
	domainerrors "github.com/michaeldcanady/go-onedrive/internal/common/errors"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal/config/domain"
	domainenv "github.com/michaeldcanady/go-onedrive/internal/core/env/domain"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal/editor/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal/profile/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockContainer struct {
	mock.Mock
}

func (m *MockContainer) Cache() domaincache.Service2 { return nil }
func (m *MockContainer) FS() domainfs.Service        { return m.Called().Get(0).(domainfs.Service) }
func (m *MockContainer) EnvironmentService() domainenv.EnvironmentService {
	return m.Called().Get(0).(domainenv.EnvironmentService)
}
func (m *MockContainer) Logger() domainlogger.LoggerService {
	return m.Called().Get(0).(domainlogger.LoggerService)
}
func (m *MockContainer) IgnoreMatcherFactory() domainfs.IgnoreMatcherFactory {
	return nil
}
func (m *MockContainer) Auth() domainauth.AuthService          { return nil }
func (m *MockContainer) Profile() domainprofile.ProfileService { return nil }
func (m *MockContainer) Config() domainconfig.ConfigService    { return nil }
func (m *MockContainer) State() domainstate.Service            { return nil }
func (m *MockContainer) Drive() domaindrive.DriveService       { return nil }
func (m *MockContainer) Account() domainaccount.Service        { return nil }
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
func (m *MockFSService) Copy(ctx context.Context, src, dst string, opts domainfs.CopyOptions) error {
	return nil
}
func (m *MockFSService) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
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
func (m *MockEnvironmentService) OutputDestination() (domainlogger.OutputDestination, error) {
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

func (m *MockLogProvider) CreateLogger(name string) (domainlogger.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(domainlogger.Logger), args.Error(1)
}
func (m *MockLogProvider) GetLogger(name string) (domainlogger.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(domainlogger.Logger), args.Error(1)
}
func (m *MockLogProvider) SetAllLevel(level string) {}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...domainlogger.Field)  { m.Called(msg, fields) }
func (m *MockLogger) Error(msg string, fields ...domainlogger.Field) { m.Called(msg, fields) }
func (m *MockLogger) Debug(msg string, fields ...domainlogger.Field) { m.Called(msg, fields) }
func (m *MockLogger) Warn(msg string, fields ...domainlogger.Field)  { m.Called(msg, fields) }
func (m *MockLogger) SetLevel(level string)                          {}
func (m *MockLogger) With(fields ...domainlogger.Field) domainlogger.Logger {
	return m
}
func (m *MockLogger) WithContext(ctx context.Context) domainlogger.Logger {
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
	args := m.Called(prefix, suffix, reader)
	var data []byte
	if args.Get(0) != nil {
		data = args.Get(0).([]byte)
	}
	return data, args.String(1), args.Error(2)
}

func (m *MockEditor) WithIO(stdin io.Reader, stdout, stderr io.Writer) domaineditor.Service {
	m.Called(stdin, stdout, stderr)
	return m
}

type MockConflictHandler struct {
	mock.Mock
}

func (m *MockConflictHandler) HandleConflict(ctx context.Context, path string, content []byte, tmpPath string) (bool, string, error) {
	args := m.Called(ctx, path, content, tmpPath)
	return args.Bool(0), args.String(1), args.Error(2)
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
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Mock Get
	mockFS.On("Get", mock.Anything, "/file.txt").Return(domainfs.Item{ETag: "tag1"}, nil)

	// Mock ReadFile
	content := "original content"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor
	mockEditor.On("LaunchTempFile", mock.Anything, ".txt", mock.Anything).
		Return(nil, "", nil)

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
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Mock Get
	mockFS.On("Get", mock.Anything, "/file.txt").Return(domainfs.Item{ETag: "tag1"}, nil)

	// Mock ReadFile
	content := "original content"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor
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
	assert.Contains(t, buf.String(), "Success:")
	assert.Contains(t, buf.String(), "updated file")
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
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Mock Get
	mockFS.On("Get", mock.Anything, "/file.txt").Return(domainfs.Item{ETag: "tag1"}, nil)

	// Mock ReadFile
	content := "original"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor: returns modified content
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
	assert.Contains(t, buf.String(), "updated file")
	mockEditor.AssertExpectations(t)
	mockFS.AssertExpectations(t)
}

func TestEditCmd_Run_ETagMismatch(t *testing.T) {
	mockContainer := new(MockContainer)
	mockFS := new(MockFSService)
	mockLogger := new(MockLogger)
	mockEditor := new(MockEditor)
	mockConflict := new(MockConflictHandler)

	mockContainer.On("FS").Return(mockFS)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Mock Get
	mockFS.On("Get", mock.Anything, "/file.txt").Return(domainfs.Item{ETag: "tag1"}, nil)

	// Mock ReadFile
	content := "original"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor: returns modified content
	mockEditor.On("LaunchTempFile", mock.Anything, ".txt", mock.Anything).
		Return([]byte("modified"), "temp-path", nil)

	// Mock WriteFile: returns ETag mismatch error
	mockFS.On("WriteFile", mock.Anything, "/file.txt", mock.Anything, mock.Anything).
		Return(domainfs.Item{}, domainerrors.ErrPrecondition)

	// Mock HandleConflict
	mockConflict.On("HandleConflict", mock.Anything, "/file.txt", []byte("modified"), "temp-path").
		Return(true, "/file.txt", nil)

	buf := new(bytes.Buffer)
	opts := Options{Path: "/file.txt", Stdout: buf, Stderr: io.Discard}

	editCmd := NewEditCmd(mockContainer).
		WithLogger(mockLogger).
		WithEditor(mockEditor).
		WithConflictHandler(mockConflict)
	err := editCmd.Run(context.Background(), opts)

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "updated file")
	mockFS.AssertExpectations(t)
	mockConflict.AssertExpectations(t)
}

func TestEditCmd_Run_SaveAsCopy(t *testing.T) {
	mockContainer := new(MockContainer)
	mockFS := new(MockFSService)
	mockLogger := new(MockLogger)
	mockEditor := new(MockEditor)
	mockConflict := new(MockConflictHandler)

	mockContainer.On("FS").Return(mockFS)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Mock Get
	mockFS.On("Get", mock.Anything, "/file.txt").Return(domainfs.Item{ETag: "tag1"}, nil)

	// Mock ReadFile
	content := "original"
	mockFS.On("ReadFile", mock.Anything, "/file.txt", mock.Anything).
		Return(io.NopCloser(strings.NewReader(content)), nil)

	// Mock Editor: returns modified content
	mockEditor.On("LaunchTempFile", mock.Anything, ".txt", mock.Anything).
		Return([]byte("modified"), "temp-path", nil)

	// Mock WriteFile: returns ETag mismatch error
	mockFS.On("WriteFile", mock.Anything, "/file.txt", mock.Anything, mock.Anything).
		Return(domainfs.Item{}, domainerrors.ErrPrecondition)

	// Mock HandleConflict: returns a different path (copy)
	mockConflict.On("HandleConflict", mock.Anything, "/file.txt", []byte("modified"), "temp-path").
		Return(true, "/file (copy).txt", nil)

	buf := new(bytes.Buffer)
	opts := Options{Path: "/file.txt", Stdout: buf, Stderr: io.Discard}

	editCmd := NewEditCmd(mockContainer).
		WithLogger(mockLogger).
		WithEditor(mockEditor).
		WithConflictHandler(mockConflict)
	err := editCmd.Run(context.Background(), opts)

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "updated file")
	assert.Contains(t, buf.String(), "/file (copy).txt")
	mockFS.AssertExpectations(t)
	mockConflict.AssertExpectations(t)
}
