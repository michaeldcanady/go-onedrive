package fs

import (
	"context"
	"io"
	"testing"

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

func TestVFS_MountAndResolve(t *testing.T) {
	v := NewVFS(nil)
	m1 := new(mockBackend)
	m2 := new(mockBackend)

	v.Mount("/od", m1)
	v.Mount("/local", m2)

	tests := []struct {
		name       string
		path       string
		wantPrefix string
		wantRel    string
		wantErr    bool
	}{
		{"exact mount", "/od", "/od", "/", false},
		{"sub path", "/od/file.txt", "/od", "/file.txt", false},
		{"nested sub path", "/od/sub/file.txt", "/od", "/sub/file.txt", false},
		{"local mount", "/local/test", "/local", "/test", false},
		{"no mount", "/unknown", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, rel, err := v.Resolve(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPrefix, prefix)
				assert.Equal(t, tt.wantRel, rel)
			}
		})
	}
}

func TestVFS_MountPrecedence(t *testing.T) {
	v := NewVFS(nil)
	m1 := new(mockBackend)
	m2 := new(mockBackend)

	v.Mount("/od", m1)
	v.Mount("/od/sub", m2)

	prefix, rel, err := v.Resolve("/od/sub/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, "/od/sub", prefix)
	assert.Equal(t, "/file.txt", rel)

	prefix, rel, err = v.Resolve("/od/other/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, "/od", prefix)
	assert.Equal(t, "/other/file.txt", rel)
}

func TestVFS_Mounts(t *testing.T) {
	v := NewVFS(nil)
	v.Mount("/od", new(mockBackend))
	v.Mount("/local", new(mockBackend))

	mounts := v.Mounts()
	assert.Len(t, mounts, 2)
	assert.Contains(t, mounts, "/od")
	assert.Contains(t, mounts, "/local")
}
