package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestVFS_Concurrency(t *testing.T) {
	v := NewVFS(nil)
	m := new(mockBackend)

	// Note: This test demonstrates that VFS is currently NOT concurrency-safe for writes (Mount)
	// but can be safely read concurrently once mounts are stable.
	// We'll test concurrent reads (Resolve) here.

	v.Mount("/od", m)

	t.Run("Concurrent Resolves", func(t *testing.T) {
		const workers = 10
		const iterations = 100

		done := make(chan bool)
		for i := 0; i < workers; i++ {
			go func() {
				for j := 0; j < iterations; j++ {
					_, _, _ = v.Resolve("/od/some/path")
				}
				done <- true
			}()
		}

		for i := 0; i < workers; i++ {
			<-done
		}
	})
}
