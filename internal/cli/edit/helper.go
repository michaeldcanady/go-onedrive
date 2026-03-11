package edit

import (
	"path/filepath"
	"strings"
)

// Name returns the filename without its extension.
func Name(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
