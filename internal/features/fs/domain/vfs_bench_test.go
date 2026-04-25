package fs

import (
	"fmt"
	"testing"
)

func BenchmarkVFS_Resolve(b *testing.B) {
	v := NewVFS(nil)
	// Setup 100 mount points to simulate a complex environment
	for i := 0; i < 100; i++ {
		v.Mount(fmt.Sprintf("/mount%d", i), new(mockBackend))
	}
	v.Mount("/mount99/sub", new(mockBackend))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = v.Resolve("/mount99/sub/file.txt")
	}
}

func BenchmarkURIFactory_FromString(b *testing.B) {
	v := NewVFS(nil)
	v.Mount("/od", new(mockBackend))
	f := NewURIFactory(v)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.FromString("/od/path/to/some/nested/file.txt")
	}
}
