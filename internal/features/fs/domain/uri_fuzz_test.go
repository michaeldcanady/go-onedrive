package fs

import (
	"testing"
)

func FuzzURIFactory_FromString(f *testing.F) {
	v := NewVFS(nil)
	v.Mount("/od", new(mockBackend))
	factory := NewURIFactory(v)

	f.Add("/od/test.txt")
	f.Add("od:/test.txt")
	f.Add("local_file.txt")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		_, _ = factory.FromString(input)
	})
}
