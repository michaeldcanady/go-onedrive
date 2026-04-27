package upload

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

func TestUpload_Functional(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		destination string
		setup       func(mRemote, mLocal *mockBackend)
		wantErr     bool
	}{
		{
			name:        "successful upload from local to remote",
			source:      "/local/src.txt",
			destination: "/od/dst.txt",
			setup: func(mRemote, mLocal *mockBackend) {
				mLocal.On("IdentityProvider").Return("")
				mLocal.On("Open", mock.Anything, mock.Anything, mock.Anything, "/src.txt").
					Return(io.NopCloser(bytes.NewBufferString("local content")), nil)

				mRemote.On("Create", mock.Anything, mock.Anything, mock.Anything, "/dst.txt", mock.Anything).
					Return(fs.Item{}, nil)
			},
			wantErr: false,
		},
		{
			name:        "upload failed - destination not found",
			source:      "/local/src.txt",
			destination: "/od/nonexistent/dst.txt",
			setup: func(mRemote, mLocal *mockBackend) {
				mLocal.On("IdentityProvider").Return("")
				mLocal.On("Open", mock.Anything, mock.Anything, mock.Anything, "/src.txt").
					Return(io.NopCloser(bytes.NewBufferString("local content")), nil)

				mRemote.On("Create", mock.Anything, mock.Anything, mock.Anything, "/nonexistent/dst.txt", mock.Anything).
					Return(fs.Item{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Setup Mocks
			mRemote := new(mockBackend)
			mLocal := new(mockBackend)
			tt.setup(mRemote, mLocal)

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
			vfs.Mount("/od", mRemote)
			vfs.Mount("/local", mLocal)

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

			mRemote.AssertExpectations(t)
			mLocal.AssertExpectations(t)
		})
	}
}
