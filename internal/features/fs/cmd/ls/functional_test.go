package ls

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
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

func TestLs_Functional(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		setup    func(m *mockBackend)
		wantErr  bool
		expected string
	}{
		{
			name: "successful ls of remote directory",
			path: "/od/test",
			setup: func(m *mockBackend) {
				m.On("IdentityProvider").Return("")
				m.On("List", mock.Anything, mock.Anything, mock.Anything, "/test").
					Return([]fs.Item{
						{Name: "file1.txt", Type: fs.TypeFile},
						{Name: "dir1", Type: fs.TypeFolder},
					}, nil)
			},
			wantErr:  false,
			expected: "file1.txt",
		},
		{
			name: "ls failed - directory not found",
			path: "/od/nonexistent",
			setup: func(m *mockBackend) {
				m.On("IdentityProvider").Return("")
				m.On("List", mock.Anything, mock.Anything, mock.Anything, "/nonexistent").
					Return(nil, errors.New("not found"))
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
			formatterFactory := formatting.NewFormatterFactory()

			handler := NewCommand(vfs, uriFactory, formatterFactory, mLog)

			// Setup Context
			buf := new(bytes.Buffer)
			cmdCtx := &CommandContext{
				Ctx: ctx,
				Options: Options{
					Path:   tt.path,
					Format: formatting.FormatShort,
					Stdout: buf,
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
				assert.Contains(t, buf.String(), tt.expected)
			}

			mBackend.AssertExpectations(t)
		})
	}
}
