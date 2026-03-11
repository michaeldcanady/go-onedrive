package app

import (
	"context"
	"io"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/fs/registry"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProvider struct {
	mock.Mock
}

func (m *mockProvider) Get(ctx context.Context, path string) (domainfs.Item, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(domainfs.Item), args.Error(1)
}

func (m *mockProvider) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	args := m.Called(ctx, path, opts)
	return args.Get(0).([]domainfs.Item), args.Error(1)
}

func (m *mockProvider) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	args := m.Called(ctx, path, opts)
	return args.Get(0).(domainfs.Item), args.Error(1)
}

func (m *mockProvider) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, path, opts)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockProvider) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	args := m.Called(ctx, path, r, opts)
	return args.Get(0).(domainfs.Item), args.Error(1)
}

func (m *mockProvider) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	args := m.Called(ctx, path, opts)
	return args.Error(0)
}

func (m *mockProvider) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	args := m.Called(ctx, path, opts)
	return args.Error(0)
}

func (m *mockProvider) Copy(ctx context.Context, src, dst string, opts domainfs.CopyOptions) error {
	args := m.Called(ctx, src, dst, opts)
	return args.Error(0)
}

func (m *mockProvider) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	args := m.Called(ctx, src, dst, opts)
	return args.Error(0)
}

func (m *mockProvider) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	args := m.Called(ctx, src, dst, opts)
	return args.Get(0).(domainfs.Item), args.Error(1)
}

func (m *mockProvider) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	args := m.Called(ctx, path, opts)
	return args.Get(0).(domainfs.Item), args.Error(1)
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		path         string
		expectedProv string
		expectedPath string
	}{
		{"onedrive:/foo/bar", "onedrive", "/foo/bar"},
		{"local://home/user", "local", "home/user"},
		{"/plain/path", "onedrive", "/plain/path"},
		{"INVALID:/path", "invalid", "/path"},
	}

	for _, tt := range tests {
		prov, subPath := util.ParsePath(tt.path)
		assert.Equal(t, tt.expectedProv, prov)
		assert.Equal(t, tt.expectedPath, subPath)
	}
}

func TestFileSystemManager_Resolve(t *testing.T) {
	reg := registry.NewRegistry()
	mockOD := new(mockProvider)
	mockLocal := new(mockProvider)

	reg.Register("onedrive", mockOD)
	reg.Register("local", mockLocal)

	fsm := NewFileSystemManager(reg)
	ctx := context.Background()

	t.Run("Explicit OneDrive", func(t *testing.T) {
		p, subPath, err := fsm.resolve(ctx, "onedrive:/foo")
		assert.NoError(t, err)
		assert.Equal(t, mockOD, p)
		assert.Equal(t, "/foo", subPath)
	})

	t.Run("Explicit Local", func(t *testing.T) {
		p, subPath, err := fsm.resolve(ctx, "local:/bar")
		assert.NoError(t, err)
		assert.Equal(t, mockLocal, p)
		assert.Equal(t, "/bar", subPath)
	})

	t.Run("Implicit OneDrive", func(t *testing.T) {
		p, subPath, err := fsm.resolve(ctx, "/baz")
		assert.NoError(t, err)
		assert.Equal(t, mockOD, p)
		assert.Equal(t, "/baz", subPath)
	})

	t.Run("Unknown prefix fallback to OneDrive", func(t *testing.T) {
		// When prefix is not registered, it's treated as full path for OneDrive
		p, subPath, err := fsm.resolve(ctx, "unknown:/path")
		assert.NoError(t, err)
		assert.Equal(t, mockOD, p)
		assert.Equal(t, "unknown:/path", subPath)
	})
}
