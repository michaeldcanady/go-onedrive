package filtering

import (
	"strings"
	"testing"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/ignore"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreFilter_IsSatisfiedBy(t *testing.T) {
	patterns := "node_modules/\n*.log\n!important.log"
	matcher, err := ignore.ParseReader(strings.NewReader(patterns))
	assert.NoError(t, err)

	f := NewIgnoreFilter(matcher)

	tests := []struct {
		name     string
		item     shared.Item
		expected bool
	}{
		{
			name: "ignored directory",
			item: shared.Item{
				Path: "node_modules",
				Type: shared.TypeFolder,
			},
			expected: false, // Should NOT satisfy (it is ignored)
		},
		{
			name: "ignored file by extension",
			item: shared.Item{
				Path: "debug.log",
				Type: shared.TypeFile,
			},
			expected: false, // Should NOT satisfy
		},
		{
			name: "negated pattern (not ignored)",
			item: shared.Item{
				Path: "important.log",
				Type: shared.TypeFile,
			},
			expected: true, // Should satisfy (not ignored)
		},
		{
			name: "regular file (not ignored)",
			item: shared.Item{
				Path: "src/main.go",
				Type: shared.TypeFile,
			},
			expected: true, // Should satisfy
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, f.IsSatisfiedBy(tt.item))
		})
	}
}
