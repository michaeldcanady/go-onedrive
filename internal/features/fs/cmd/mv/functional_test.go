package mv

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBackend struct {
	mock.Mock
}

func (m *mockBackend) Name() string { return m.Called().String(0) }
func (m *mockBackend) IdentityProvider() string { return m.Called().String(0) }
func (m *mockBackend) Stat(ctx context.Context, token, driveID, path string) (fs.Item, error) {
	args := m.Called(ctx, token, driveID, path)
	return args.Get(0).(fs.Item), args.Error(1)
}
func (m *mockBackend) List(ctx context.Context, token, driveID, path string) ([]fs.Item, error) {
	args := m.Called(ctx, token, driveID, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]fs.Item), args.Error(1)
}
func (m *mockBackend) Open(ctx context.Context, token, driveID, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, token, driveID, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
func (m *mockBackend) Create(ctx context.Context, token, driveID, path string, r io.Reader) (fs.Item, error) {
	args := m.Called(ctx, token, driveID, path, r)
	return args.Get(0).(fs.Item), args.Error(1)
}
func (m *mockBackend) Mkdir(ctx context.Context, token, driveID, path string) error {
	return m.Called(ctx, token, driveID, path).Error(0)
}
func (m *mockBackend) Remove(ctx context.Context, token, driveID, path string) error {
	return m.Called(ctx, token, driveID, path).Error(0)
}
func (m *mockBackend) Capabilities() fs.Capabilities {
	return m.Called().Get(0).(fs.Capabilities)
}
func (m *mockBackend) Move(ctx context.Context, token, driveID, src, dst string) error {
	return m.Called(ctx, token, driveID, src, dst).Error(0)
}
func (m *mockBackend) Copy(ctx context.Context, token, driveID, src, dst string) error {
	return m.Called(ctx, token, driveID, src, dst).Error(0)
}

type mockLogger struct{ mock.Mock }
func (m *mockLogger) Debug(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) Error(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) Warn(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *mockLogger) Info(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *mockLogger) SetLevel(level logger.Level)              { m.Called(level) }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger { return m.Called(ctx).Get(0).(logger.Logger) }
func (m *mockLogger) With(fields ...logger.Field) logger.Logger { return m.Called(fields).Get(0).(logger.Logger) }

func TestMv_Functional(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		destination string
		setup       func(m *mockBackend)
		wantErr     bool
	}{
		{
			name:        "successful move (native)",
			source:      "/od/src.txt",
			destination: "/od/dst.txt",
			setup: func(m *mockBackend) {
				m.On("IdentityProvider").Return("")
				m.On("Capabilities").Return(fs.Capabilities{CanMove: true})
				m.On("Move", mock.Anything, mock.Anything, mock.Anything, "/src.txt", "/dst.txt").Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "move failed - source not found",
			source:      "/od/nonexistent.txt",
			destination: "/od/dst.txt",
			setup: func(m *mockBackend) {
				m.On("IdentityProvider").Return("")
				m.On("Capabilities").Return(fs.Capabilities{CanMove: true})
				m.On("Move", mock.Anything, mock.Anything, mock.Anything, "/nonexistent.txt", "/dst.txt").
					Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Setup Mocks
			mBackend := new(mockBackend)
			tt.setup(mBackend)

			mLog := new(mockLogger)
			mLog.On("WithContext", mock.Anything).Return(mLog)
			mLog.On("With", mock.Anything).Return(mLog)
			mLog.On("Info", mock.Anything, mock.Anything).Return()
			mLog.On("Debug", mock.Anything, mock.Anything).Return()
			if tt.wantErr {
				mLog.On("Error", mock.Anything, mock.Anything).Return()
			}

			// Setup Real Services
			vfs := fsdomain.NewVFS(nil)
			vfs.Mount("/od", mBackend)

			uriFactory := fsdomain.NewURIFactory(vfs)

			handler := NewCommand(vfs, uriFactory, mLog)

			// Setup Context
			buf := new(bytes.Buffer)
			cmdCtx := &CommandContext{
				Ctx: ctx,
				Options: Options{
					Source:      tt.source,
					Destination: tt.destination,
					Stdout:      buf,
				},
			}

			// Execute
			err := handler.Validate(cmdCtx)
			assert.NoError(t, err)

			err = handler.Execute(cmdCtx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mBackend.AssertExpectations(t)
		})
	}
}
