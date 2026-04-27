package mkdir

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBackend struct {
	mock.Mock
}

func (m *mockBackend) Name() string             { return m.Called().String(0) }
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

func TestMkdir_Functional(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		setup   func(m *mockBackend)
		wantErr bool
	}{
		{
			name: "successful mkdir",
			path: "/od/new-dir",
			setup: func(m *mockBackend) {
				m.On("IdentityProvider").Return("")
				m.On("Mkdir", mock.Anything, mock.Anything, mock.Anything, "/new-dir").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "mkdir failed - already exists",
			path: "/od/existing-dir",
			setup: func(m *mockBackend) {
				m.On("IdentityProvider").Return("")
				m.On("Mkdir", mock.Anything, mock.Anything, mock.Anything, "/existing-dir").
					Return(errors.New("already exists"))
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
					Path:   tt.path,
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
			}

			mBackend.AssertExpectations(t)
		})
	}
}
