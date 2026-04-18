package filtering

import (
	"testing"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/stretchr/testify/assert"
)

func TestItemTypeFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   shared.ItemType
		item     shared.Item
		expected bool
	}{
		{
			name:   "filter for files: item is file",
			filter: shared.TypeFile,
			item: shared.Item{
				Name: "file.txt",
				Type: shared.TypeFile,
			},
			expected: true,
		},
		{
			name:   "filter for files: item is folder",
			filter: shared.TypeFile,
			item: shared.Item{
				Name: "folder",
				Type: shared.TypeFolder,
			},
			expected: false,
		},
		{
			name:   "filter for folders: item is folder",
			filter: shared.TypeFolder,
			item: shared.Item{
				Name: "folder",
				Type: shared.TypeFolder,
			},
			expected: true,
		},
		{
			name:   "filter for folders: item is file",
			filter: shared.TypeFolder,
			item: shared.Item{
				Name: "file.txt",
				Type: shared.TypeFile,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ItemTypeFilter{itemType: tt.filter}
			assert.Equal(t, tt.expected, i.IsSatisfiedBy(tt.item))
		})
	}
}
