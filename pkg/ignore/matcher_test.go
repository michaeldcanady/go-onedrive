package ignore

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	input := `
node_modules/
*.log
!keep.log
temp-*
/absolute-dir/
`
	m, err := ParseReader(strings.NewReader(input))
	assert.NoError(t, err)

	tests := []struct {
		path   string
		isDir  bool
		ignore bool
	}{
		{path: "node_modules", isDir: true, ignore: true},
		{path: "node_modules/foo", isDir: false, ignore: true},
		{path: "error.log", isDir: false, ignore: true},
		{path: "keep.log", isDir: false, ignore: false},
		{path: "src/main.go", isDir: false, ignore: false},
		{path: "temp-cache", isDir: true, ignore: true},
		{path: "temp-files/data.txt", isDir: false, ignore: true},
		{path: "absolute-dir/file", isDir: false, ignore: true},
		{path: "./node_modules", isDir: true, ignore: true},
		{path: "/node_modules", isDir: true, ignore: true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.ignore, m.ShouldIgnore(tt.path, tt.isDir))
		})
	}
}

func TestNewMatcher(t *testing.T) {
	patterns := []Pattern{
		{Path: "ignore-me", IsNegate: false, IsDir: false},
	}
	m := NewMatcher(patterns)
	assert.NotNil(t, m)
	assert.True(t, m.ShouldIgnore("ignore-me", false))
}

func TestMatcher_EmptyPatterns(t *testing.T) {
	m := NewMatcher(nil)
	assert.False(t, m.ShouldIgnore("anything", false))
}

func TestMatcher_NegationPrecedence(t *testing.T) {
	input := `
*.txt
!important.txt
`
	m, err := ParseReader(strings.NewReader(input))
	assert.NoError(t, err)

	assert.True(t, m.ShouldIgnore("file.txt", false))
	assert.False(t, m.ShouldIgnore("important.txt", false))
}

func TestMatcher_DirectoryOnly(t *testing.T) {
	input := `
build/
`
	m, err := ParseReader(strings.NewReader(input))
	assert.NoError(t, err)

	assert.True(t, m.ShouldIgnore("build", true))
	assert.False(t, m.ShouldIgnore("build", false))
}

func TestMatcher_LastPatternWins(t *testing.T) {
	input := `
*.txt
!important.txt
ignored.txt
`
	m, err := ParseReader(strings.NewReader(input))
	assert.NoError(t, err)

	assert.True(t, m.ShouldIgnore("ignored.txt", false))
}

func TestMatcher_MultipleSegments(t *testing.T) {
	input := `
a/b
`
	m, err := ParseReader(strings.NewReader(input))
	assert.NoError(t, err)

	assert.True(t, m.ShouldIgnore("a/b/c", false))
	assert.False(t, m.ShouldIgnore("a/x/b", false))
}
