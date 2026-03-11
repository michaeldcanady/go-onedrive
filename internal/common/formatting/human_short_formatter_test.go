package formatting

import (
	"bytes"
	"testing"

	fs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
	"github.com/stretchr/testify/assert"
)

type mockTerminalInfo struct {
	width int
}

func (m mockTerminalInfo) Width() int {
	return m.width
}

func TestHumanShortFormatter_Format(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		items    []fs.Item
		expected []string // substrings that should be present
	}{
		{
			name:     "Single item",
			width:    80,
			items:    []fs.Item{{Name: "file1"}},
			expected: []string{"file1"},
		},
		{
			name:     "Multiple items, narrow terminal",
			width:    10,
			items:    []fs.Item{{Name: "file1"}, {Name: "file2"}},
			expected: []string{"file1", "file2"},
		},
		{
			name:     "Multiple items, wide terminal",
			width:    80,
			items:    []fs.Item{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &HumanShortFormatter{term: mockTerminalInfo{width: tt.width}}
			buf := new(bytes.Buffer)
			err := f.Format(buf, tt.items)
			assert.NoError(t, err)

			output := buf.String()
			for _, exp := range tt.expected {
				assert.Contains(t, output, exp)
			}
		})
	}
}
