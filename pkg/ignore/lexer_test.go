package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Basic pattern",
			input: "file.txt",
			expected: []Token{
				{TokenText, "file.txt"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Comment and negation",
			input: "# comment\n!important",
			expected: []Token{
				{TokenHash, "#"},
				{TokenSpace, " "},
				{TokenText, "comment"},
				{TokenNewline, "\n"},
				{TokenBang, "!"},
				{TokenText, "important"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Directory and spaces",
			input: "/dir/ ",
			expected: []Token{
				{TokenSlash, "/"},
				{TokenText, "dir"},
				{TokenSlash, "/"},
				{TokenSpace, " "},
				{TokenEOF, ""},
			},
		},
		{
			name:  "CRLF handling",
			input: "file\r\nnext",
			expected: []Token{
				{TokenText, "file"},
				{TokenNewline, "\n"},
				{TokenText, "next"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Empty input",
			input: "",
			expected: []Token{
				{TokenEOF, ""},
			},
		},
		{
			name:  "Only newline \\r",
			input: "\r",
			expected: []Token{
				{TokenNewline, "\n"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Multiple spaces",
			input: "  ",
			expected: []Token{
				{TokenSpace, " "},
				{TokenSpace, " "},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Special characters in text",
			input: "file*?[]!",
			expected: []Token{
				{TokenText, "file*?[]"},
				{TokenBang, "!"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Escaped special character",
			input: "\\#file",
			expected: []Token{
				{TokenText, "#file"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Escaped exclamation",
			input: "\\!important",
			expected: []Token{
				{TokenText, "!important"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Escaped space",
			input: "file\\ name",
			expected: []Token{
				{TokenText, "file"},
				{TokenText, " name"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "Unexpected EOF after escape",
			input: "file\\",
			expected: []Token{
				{TokenText, "file"},
				{TokenError, "unexpected EOF after escape"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			for _, exp := range tt.expected {
				token := l.NextToken()
				assert.Equal(t, exp.Type, token.Type)
				if exp.Type != TokenEOF {
					assert.Equal(t, exp.Literal, token.Literal)
				}
			}
		})
	}
}
