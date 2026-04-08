package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Pattern
	}{
		{
			name: "Various patterns",
			input: `
# a comment
!include.me
exclude/
normal-file
  spaced-file  
`,
			expected: []Pattern{
				{Original: "!include.me", Path: "include.me", IsNegate: true, IsDir: false},
				{Original: "exclude/", Path: "exclude", IsNegate: false, IsDir: true},
				{Original: "normal-file", Path: "normal-file", IsNegate: false, IsDir: false},
				{Original: "spaced-file", Path: "spaced-file", IsNegate: false, IsDir: false},
			},
		},
		{
			name:     "Empty and whitespace lines",
			input:    "\n  \n\t\n",
			expected: nil,
		},
		{
			name:  "Inline comments",
			input: "file.txt # keep this",
			expected: []Pattern{
				{Original: "file.txt", Path: "file.txt", IsNegate: false, IsDir: false},
			},
		},
		{
			name:     "Only comments",
			input:    "# comment1\n# comment2",
			expected: nil,
		},
		{
			name:  "Empty file",
			input: "",
			expected: nil,
		},
		{
			name:  "Multiple negations",
			input: "!!important.txt",
			expected: []Pattern{
				{Original: "!!important.txt", Path: "!important.txt", IsNegate: true, IsDir: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			p := NewParser(l)
			patterns := p.Parse()

			if tt.expected == nil {
				assert.Empty(t, patterns)
			} else {
				assert.Equal(t, len(tt.expected), len(patterns))
				for i, exp := range tt.expected {
					assert.Equal(t, exp.Path, patterns[i].Path, "Path mismatch at index %d", i)
					assert.Equal(t, exp.IsNegate, patterns[i].IsNegate, "IsNegate mismatch at index %d", i)
					assert.Equal(t, exp.IsDir, patterns[i].IsDir, "IsDir mismatch at index %d", i)
				}
			}
		})
	}
}
