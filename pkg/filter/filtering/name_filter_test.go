package filtering

import (
	"testing"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/stretchr/testify/assert"
)

func TestNameFilter(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		item     shared.Item
		expected bool
	}{
		{
			name:     "no patterns: matches anything",
			patterns: []string{},
			item: shared.Item{
				Name: "file.txt",
			},
			expected: true,
		},
		{
			name:     "single pattern: match",
			patterns: []string{"*.txt"},
			item: shared.Item{
				Name: "file.txt",
			},
			expected: true,
		},
		{
			name:     "single pattern: mismatch",
			patterns: []string{"*.jpg"},
			item: shared.Item{
				Name: "file.txt",
			},
			expected: false,
		},
		{
			name:     "multiple patterns: match first",
			patterns: []string{"*.txt", "*.jpg"},
			item: shared.Item{
				Name: "file.txt",
			},
			expected: true,
		},
		{
			name:     "multiple patterns: match second",
			patterns: []string{"*.jpg", "*.txt"},
			item: shared.Item{
				Name: "file.txt",
			},
			expected: true,
		},
		{
			name:     "multiple patterns: mismatch all",
			patterns: []string{"*.jpg", "*.png"},
			item: shared.Item{
				Name: "file.txt",
			},
			expected: false,
		},
		{
			name:     "exact match",
			patterns: []string{"readme.md"},
			item: shared.Item{
				Name: "readme.md",
			},
			expected: true,
		},
		{
			name:     "invalid pattern: handled gracefully",
			patterns: []string{"["}, // Invalid glob
			item: shared.Item{
				Name: "file.txt",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNameFilter(tt.patterns)
			assert.Equal(t, tt.expected, n.IsSatisfiedBy(tt.item))
		})
	}
}
