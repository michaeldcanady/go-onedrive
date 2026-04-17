package fs

import (
	"context"
	"io"
	"testing"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
)

type mockBackend struct {
	name string
}

func (m *mockBackend) Name() string { return m.name }
func (m *mockBackend) Stat(ctx context.Context, path string) (fs.Item, error) {
	return fs.Item{Path: path, Name: m.name}, nil
}
func (m *mockBackend) List(ctx context.Context, path string) ([]fs.Item, error) { return nil, nil }
func (m *mockBackend) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return nil, nil
}
func (m *mockBackend) Create(ctx context.Context, path string, r io.Reader) (fs.Item, error) {
	return fs.Item{}, nil
}
func (m *mockBackend) Mkdir(ctx context.Context, path string) error  { return nil }
func (m *mockBackend) Remove(ctx context.Context, path string) error { return nil }
func (m *mockBackend) Capabilities() fs.Capabilities               { return fs.Capabilities{} }

func TestVFS_Resolve(t *testing.T) {
	vfs := NewVFS()
	local := &mockBackend{name: "local"}
	work := &mockBackend{name: "work"}
	personal := &mockBackend{name: "personal"}

	vfs.Mount("/", local)
	vfs.Mount("/work", work)
	vfs.Mount("/personal/onedrive", personal)

	tests := []struct {
		path         string
		expectedName string
		expectedRel  string
	}{
		{"/file.txt", "local", "/file.txt"},
		{"/work/project/main.go", "work", "/project/main.go"},
		{"/work", "work", "/"},
		{"/personal/onedrive/photos/me.jpg", "personal", "/photos/me.jpg"},
		{"/other/path", "local", "/other/path"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			backend, rel, err := vfs.resolve(tt.path)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedName, backend.Name())
			assert.Equal(t, tt.expectedRel, rel)
		})
	}
}

func TestVFS_SelectBackend(t *testing.T) {
	vfs := NewVFS()
	local := &mockBackend{name: "local"}
	work := &mockBackend{name: "work"}

	vfs.Mount("/local", local)
	vfs.Mount("/work", work)

	t.Run("by provider", func(t *testing.T) {
		uri := &fs.URI{Provider: "/work", Path: "/some/path"}
		backend, rel, err := vfs.selectBackend(uri)
		assert.NoError(t, err)
		assert.Equal(t, "work", backend.Name())
		assert.Equal(t, "/some/path", rel)
	})

	t.Run("by path fallback", func(t *testing.T) {
		uri := &fs.URI{Path: "/local/file.txt"}
		backend, rel, err := vfs.selectBackend(uri)
		assert.NoError(t, err)
		assert.Equal(t, "local", backend.Name())
		assert.Equal(t, "/file.txt", rel)
	})

	t.Run("provider not found fallback to path", func(t *testing.T) {
		uri := &fs.URI{Provider: "/missing", Path: "/work/file.txt"}
		backend, rel, err := vfs.selectBackend(uri)
		assert.NoError(t, err)
		assert.Equal(t, "work", backend.Name())
		assert.Equal(t, "/file.txt", rel)
	})
}
