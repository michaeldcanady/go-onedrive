package filtering

import (
	"testing"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/stretchr/testify/assert"
)

func TestSizeFilter(t *testing.T) {
	tests := []struct {
		name     string
		minSize  *int64
		maxSize  *int64
		item     shared.Item
		expected bool
	}{
		{
			name:     "no bounds: matches anything",
			minSize:  nil,
			maxSize:  nil,
			item:     shared.Item{Size: 100},
			expected: true,
		},
		{
			name:     "min size: below threshold",
			minSize:  ptr(int64(100)),
			maxSize:  nil,
			item:     shared.Item{Size: 50},
			expected: false,
		},
		{
			name:     "min size: at threshold",
			minSize:  ptr(int64(100)),
			maxSize:  nil,
			item:     shared.Item{Size: 100},
			expected: true,
		},
		{
			name:     "min size: above threshold",
			minSize:  ptr(int64(100)),
			maxSize:  nil,
			item:     shared.Item{Size: 150},
			expected: true,
		},
		{
			name:     "max size: above threshold",
			minSize:  nil,
			maxSize:  ptr(int64(200)),
			item:     shared.Item{Size: 250},
			expected: false,
		},
		{
			name:     "max size: at threshold",
			minSize:  nil,
			maxSize:  ptr(int64(200)),
			item:     shared.Item{Size: 200},
			expected: true,
		},
		{
			name:     "max size: below threshold",
			minSize:  nil,
			maxSize:  ptr(int64(200)),
			item:     shared.Item{Size: 150},
			expected: true,
		},
		{
			name:     "both bounds: below range",
			minSize:  ptr(int64(100)),
			maxSize:  ptr(int64(200)),
			item:     shared.Item{Size: 50},
			expected: false,
		},
		{
			name:     "both bounds: above range",
			minSize:  ptr(int64(100)),
			maxSize:  ptr(int64(200)),
			item:     shared.Item{Size: 250},
			expected: false,
		},
		{
			name:     "both bounds: within range",
			minSize:  ptr(int64(100)),
			maxSize:  ptr(int64(200)),
			item:     shared.Item{Size: 150},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSizeFilter(tt.minSize, tt.maxSize)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy(tt.item))
		})
	}
}
