package filtering

import (
	"testing"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/stretchr/testify/assert"
)

func TestHiddenFilter(t *testing.T) {
	tests := []struct {
		name     string
		hidden   bool
		item     shared.Item
		expected bool
	}{
		{
			name:   "exclude hidden: normal file",
			hidden: false,
			item: shared.Item{
				Name: "file.txt",
			},
			expected: true,
		},
		{
			name:   "exclude hidden: hidden file",
			hidden: false,
			item: shared.Item{
				Name: ".hidden",
			},
			expected: false,
		},
		{
			name:   "include only hidden: normal file",
			hidden: true,
			item: shared.Item{
				Name: "file.txt",
			},
			expected: false,
		},
		{
			name:   "include only hidden: hidden file",
			hidden: true,
			item: shared.Item{
				Name: ".hidden",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HiddenFilter{hidden: tt.hidden}
			assert.Equal(t, tt.expected, h.IsSatisfiedBy(tt.item))
		})
	}
}
