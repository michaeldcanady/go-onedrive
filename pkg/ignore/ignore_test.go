package ignore

import (
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	input := `# comment
!ignore
/dir/`
	l := NewLexer(input)

	tokens := []Token{
		{TokenHash, "#"},
		{TokenSpace, " "},
		{TokenText, "comment"},
		{TokenNewline, "\n"},
		{TokenBang, "!"},
		{TokenText, "ignore"},
		{TokenNewline, "\n"},
		{TokenSlash, "/"},
		{TokenText, "dir"},
		{TokenSlash, "/"},
		{TokenEOF, ""},
	}

	for _, expected := range tokens {
		token := l.NextToken()
		assert.Equal(t, expected.Type, token.Type)
		if expected.Type != TokenEOF {
			assert.Equal(t, expected.Literal, token.Literal)
		}
	}
}

func TestParser(t *testing.T) {
	input := `
# a comment
!include.me
exclude/
normal-file
`
	l := NewLexer(input)
	p := NewParser(l)
	patterns := p.Parse()

	expected := []Pattern{
		{Original: "!include.me", Path: "include.me", IsNegate: true, IsDir: false},
		{Original: "exclude/", Path: "exclude", IsNegate: false, IsDir: true},
		{Original: "normal-file", Path: "normal-file", IsNegate: false, IsDir: false},
	}

	assert.Equal(t, len(expected), len(patterns))
	for i, exp := range expected {
		assert.Equal(t, exp.Path, patterns[i].Path)
		assert.Equal(t, exp.IsNegate, patterns[i].IsNegate)
		assert.Equal(t, exp.IsDir, patterns[i].IsDir)
	}
}

func TestMatcher(t *testing.T) {
	input := `
node_modules/
*.log
!keep.log
`
	m, err := ParseReader(strings.NewReader(input))
	assert.NoError(t, err)

	tests := []struct {
		path   string
		isDir  bool
		ignore bool
	}{
		{"node_modules", true, true},
		{"node_modules/foo", false, true},
		{"error.log", false, true},
		{"keep.log", false, false},
		{"src/main.go", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.ignore, m.ShouldIgnore(tt.path, tt.isDir))
		})
	}
}
