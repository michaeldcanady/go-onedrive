package filtering

import (
	"testing"
	"time"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/stretchr/testify/assert"
)

func TestOptionsFilterer_Filter(t *testing.T) {
	now := time.Now()
	items := []shared.Item{
		{Name: "file1.txt", Type: shared.TypeFile, Size: 100, ModifiedAt: now},
		{Name: "file2.jpg", Type: shared.TypeFile, Size: 200, ModifiedAt: now.Add(-time.Hour)},
		{Name: ".hidden", Type: shared.TypeFile, Size: 50, ModifiedAt: now},
		{Name: "folder1", Type: shared.TypeFolder, Size: 0, ModifiedAt: now},
	}

	tests := []struct {
		name     string
		opts     FilterOptions
		expected []string
	}{
		{
			name:     "default options (exclude hidden)",
			opts:     *NewFilterOptions(),
			expected: []string{"file1.txt", "file2.jpg", "folder1"},
		},
		{
			name: "include all",
			opts: FilterOptions{
				IncludeAll: true,
			},
			expected: []string{"file1.txt", "file2.jpg", ".hidden", "folder1"},
		},
		{
			name: "filter by type (files only)",
			opts: FilterOptions{
				ItemType:   shared.TypeFile,
				IncludeAll: true,
			},
			expected: []string{"file1.txt", "file2.jpg", ".hidden"},
		},
		{
			name: "filter by name glob",
			opts: FilterOptions{
				Names:      []string{"*.txt"},
				IncludeAll: true,
			},
			expected: []string{"file1.txt"},
		},
		{
			name: "filter by multiple name globs",
			opts: FilterOptions{
				Names:      []string{"*.txt", "*.jpg"},
				IncludeAll: true,
			},
			expected: []string{"file1.txt", "file2.jpg"},
		},
		{
			name: "filter by min size",
			opts: FilterOptions{
				MinSize:    ptr(int64(150)),
				IncludeAll: true,
			},
			expected: []string{"file2.jpg"},
		},
		{
			name: "filter by max size",
			opts: FilterOptions{
				MaxSize:    ptr(int64(150)),
				IncludeAll: true,
			},
			expected: []string{"file1.txt", ".hidden", "folder1"},
		},
		{
			name: "filter by date (modified after)",
			opts: FilterOptions{
				ModifiedAfter: ptr(now.Add(-30 * time.Minute)),
				IncludeAll:    true,
			},
			expected: []string{"file1.txt", ".hidden", "folder1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewOptionsFilterer(tt.opts)
			filtered, err := f.Filter(items)
			assert.NoError(t, err)

			var names []string
			for _, item := range filtered {
				names = append(names, item.Name)
			}
			assert.ElementsMatch(t, tt.expected, names)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
