package cp

import (
	"bytes"
	"context"
	"io"
	"testing"

	environment "github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	drive "github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	editor "github.com/michaeldcanady/go-onedrive/internal/features/editor/domain"
	registry "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockContainer struct{ mock.Mock }

func (m *mockContainer) Logger() logger.Service     { return m.Called().Get(0).(logger.Service) }
func (m *mockContainer) Config() config.Service     { return m.Called().Get(0).(config.Service) }
func (m *mockContainer) Mounts() mount.Service      { return m.Called().Get(0).(mount.Service) }
func (m *mockContainer) Identity() identity.Service { return m.Called().Get(0).(identity.Service) }
func (m *mockContainer) Profile() profile.Service   { return m.Called().Get(0).(profile.Service) }
func (m *mockContainer) FS() registry.Service       { return m.Called().Get(0).(registry.Service) }
func (m *mockContainer) Environment() environment.Service {
	return m.Called().Get(0).(environment.Service)
}
func (m *mockContainer) Editor() editor.Service { return m.Called().Get(0).(editor.Service) }
func (m *mockContainer) Drive() drive.Service   { return m.Called().Get(0).(drive.Service) }
func (m *mockContainer) URIFactory() *registry.URIFactory {
	return m.Called().Get(0).(*registry.URIFactory)
}

type mockLoggerService struct{ mock.Mock }

func (m *mockLoggerService) CreateLogger(name string) (logger.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(logger.Logger), args.Error(1)
}
func (m *mockLoggerService) SetAllLevel(level logger.Level) { m.Called(level) }
func (m *mockLoggerService) Reconfigure(level logger.Level, output string, format string) error {
	return m.Called(level, output, format).Error(0)
}

type mockLogger struct{ mock.Mock }

func (m *mockLogger) Debug(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) Error(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) Warn(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *mockLogger) SetLevel(level logger.Level)              { m.Called(level) }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger {
	return m.Called(ctx).Get(0).(logger.Logger)
}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger {
	return m.Called(fields).Get(0).(logger.Logger)
}
func (m *mockLogger) Info(msg string, fields ...logger.Field) { m.Called(msg, fields) }

type mockFsService struct{ mock.Mock }

func (m *mockFsService) Name() string { return m.Called().String(0) }
func (m *mockFsService) Get(ctx context.Context, uri *pkgfs.URI) (pkgfs.Item, error) {
	args := m.Called(ctx, uri)
	return args.Get(0).(pkgfs.Item), args.Error(1)
}
func (m *mockFsService) List(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ListOptions) ([]pkgfs.Item, error) {
	args := m.Called(ctx, uri, opts)
	return args.Get(0).([]pkgfs.Item), args.Error(1)
}
func (m *mockFsService) ReadFile(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ReadOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, uri, opts)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
func (m *mockFsService) Stat(ctx context.Context, uri *pkgfs.URI) (pkgfs.Item, error) {
	args := m.Called(ctx, uri)
	return args.Get(0).(pkgfs.Item), args.Error(1)
}
func (m *mockFsService) WriteFile(ctx context.Context, uri *pkgfs.URI, r io.Reader, opts pkgfs.WriteOptions) (pkgfs.Item, error) {
	args := m.Called(ctx, uri, r, opts)
	return args.Get(0).(pkgfs.Item), args.Error(1)
}
func (m *mockFsService) Mkdir(ctx context.Context, uri *pkgfs.URI) error {
	return m.Called(ctx, uri).Error(0)
}
func (m *mockFsService) Touch(ctx context.Context, uri *pkgfs.URI) (pkgfs.Item, error) {
	args := m.Called(ctx, uri)
	return args.Get(0).(pkgfs.Item), args.Error(1)
}
func (m *mockFsService) Remove(ctx context.Context, uri *pkgfs.URI) error {
	return m.Called(ctx, uri).Error(0)
}
func (m *mockFsService) Copy(ctx context.Context, src, dst *pkgfs.URI, opts pkgfs.CopyOptions) error {
	return m.Called(ctx, src, dst, opts).Error(0)
}
func (m *mockFsService) Move(ctx context.Context, src, dst *pkgfs.URI) error {
	return m.Called(ctx, src, dst).Error(0)
}

type mockVFS struct{ mock.Mock }

func (m *mockVFS) Resolve(absPath string) (string, string, error) {
	args := m.Called(absPath)
	return args.String(0), args.String(1), args.Error(2)
}

func TestCpCommand_Integration(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		setup   func(m *mockContainer, mFS *mockFsService, mLogSvc *mockLoggerService, mLog *mockLogger, mVFS *mockVFS)
		wantErr bool
	}{
		{
			name: "cp successfully",
			args: []string{"cp", "od:/src", "od:/dst"},
			setup: func(m *mockContainer, mFS *mockFsService, mLogSvc *mockLoggerService, mLog *mockLogger, mVFS *mockVFS) {
				mLogSvc.On("CreateLogger", "cp").Return(mLog, nil)
				m.On("Logger").Return(mLogSvc)
				m.On("FS").Return(mFS)
				m.On("URIFactory").Return(registry.NewURIFactory(mVFS))
				mVFS.On("Resolve", mock.Anything).Return("od:/", "/", nil)
				mLog.On("WithContext", mock.Anything).Return(mLog)
				mLog.On("With", mock.Anything).Return(mLog)
				mLog.On("Info", mock.Anything, mock.Anything).Return()
				mLog.On("Debug", mock.Anything, mock.Anything).Return()
				mFS.On("Copy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mContainer := new(mockContainer)
			mFS := new(mockFsService)
			mLogSvc := new(mockLoggerService)
			mLog := new(mockLogger)
			mVFS := new(mockVFS)

			if tt.setup != nil {
				tt.setup(mContainer, mFS, mLogSvc, mLog, mVFS)
			}

			cmd := CreateCpCmd(mContainer)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args[1:])

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
