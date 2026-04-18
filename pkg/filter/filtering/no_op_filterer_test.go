package filtering

import (
	"testing"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/stretchr/testify/assert"
)

func TestNoOpFilterer_Filter(t *testing.T) {
	tests := []struct {
		name     string
		items    []shared.Item
		expected []shared.Item
	}{
		{
			name:     "empty slice",
			items:    []shared.Item{},
			expected: []shared.Item{},
		},
		{
			name: "single item",
			items: []shared.Item{
				{Name: "file1.txt"},
			},
			expected: []shared.Item{
				{Name: "file1.txt"},
			},
		},
		{
			name: "multiple items",
			items: []shared.Item{
				{Name: "file1.txt"},
				{Name: "folder1"},
				{Name: ".hidden"},
			},
			expected: []shared.Item{
				{Name: "file1.txt"},
				{Name: "folder1"},
				{Name: ".hidden"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewNoOpFilterer()
			filtered, err := f.Filter(tt.items)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, filtered)

			// Ensure it's a copy, not the same slice reference (though items themselves are shallow copied)
			if len(tt.items) > 0 {
				assert.NotSame(t, &tt.items[0], &filtered[0])
			}
		})
	}
}
