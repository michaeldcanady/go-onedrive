package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	input := `# This is a comment
!negated
dir/
foo.txt
nested/**/bar`
	l := NewLexer(input)

	var tokens []Token
	for {
		tok := l.NextToken()
		if tok.Type == TokenEOF {
			break
		}
		tokens = append(tokens, tok)
	}

	expectedTypes := []TokenType{
		TokenComment, TokenText, TokenNewline, // # This is a comment

		TokenNegate, TokenText, TokenNewline, // !negated

		TokenText, TokenSlash, TokenNewline, // dir/

		TokenText, TokenNewline, // foo.txt

		TokenText, TokenSlash, TokenDoubleStar, TokenSlash, TokenText, // nested/**/bar
	}

	assert.Equal(t, len(expectedTypes), len(tokens), "Token count mismatch")

	for i, typ := range expectedTypes {
		if i < len(tokens) {
			assert.Equal(t, typ, tokens[i].Type, "Token %d type mismatch", i)
		}
	}
}

func TestParser(t *testing.T) {
	input := `
# Comment
!ignore.me
/root
foo/
bar/**/baz
escaped\#hash
trailing   `
	l := NewLexer(input)
	p := NewParser(l)
	doc := p.Parse()

	expectedRules := []Rule{
		{Pattern: "ignore.me", Negate: true, DirOnly: false},
		{Pattern: "/root", Negate: false, DirOnly: false},
		{Pattern: "foo", Negate: false, DirOnly: true},
		{Pattern: "bar/**/baz", Negate: false, DirOnly: false},
		{Pattern: "escaped#hash", Negate: false, DirOnly: false},
		{Pattern: "trailing", Negate: false, DirOnly: false}, // Space ignored
	}

	assert.Equal(t, len(expectedRules), len(doc.Rules), "Rule count mismatch")

	for i, exp := range expectedRules {
		if i < len(doc.Rules) {
			assert.Equal(t, exp.Pattern, doc.Rules[i].Pattern, "Rule %d pattern mismatch", i)
			assert.Equal(t, exp.Negate, doc.Rules[i].Negate, "Rule %d negate mismatch", i)
			assert.Equal(t, exp.DirOnly, doc.Rules[i].DirOnly, "Rule %d dirOnly mismatch", i)
		}
	}
}

func TestMatcher(t *testing.T) {
	rules := `
node_modules/
*.log
!important.log
/src/
doc/*.md
`
	l := NewLexer(rules)
	p := NewParser(l)
	doc := p.Parse()
	m := NewMatcher(doc.Rules)

	tests := []struct {
		path   string
		isDir  bool
		ignore bool
	}{
		{"node_modules", true, true},
		{"node_modules/foo.js", false, true},
		{"src", true, true},
		{"src/main.go", false, true}, // Implicitly ignored because parent is ignored
		{"foo/src", true, false},     // Leading slash anchors to root
		{"error.log", false, true},
		{"important.log", false, false},
		{"doc/readme.md", false, true},
		{"doc/other/readme.md", false, false}, // * matches no slash
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			ignored := m.Match(tt.path, tt.isDir)
			assert.Equal(t, tt.ignore, ignored, "Path: %s", tt.path)
		})
	}
}

// Additional test for escaped spaces and special characters
func TestParserSpecial(t *testing.T) {
	input := `
space\ 
\#hash
\!bang
`
	l := NewLexer(input)
	p := NewParser(l)
	doc := p.Parse()

	expectedRules := []Rule{
		{Pattern: "space ", Negate: false, DirOnly: false},
		{Pattern: "#hash", Negate: false, DirOnly: false},
		{Pattern: "!bang", Negate: false, DirOnly: false},
	}

	assert.Equal(t, len(expectedRules), len(doc.Rules))
	for i, exp := range expectedRules {
		assert.Equal(t, exp.Pattern, doc.Rules[i].Pattern)
	}
}
