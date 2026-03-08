package ignore

// TokenType defines the type of a lexical token.
type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenText
	TokenSlash
	TokenBang
	TokenHash
	TokenNewline
	TokenSpace
)

const (
	charHash    = '#'
	charBang    = '!'
	charSlash   = '/'
	charSpace   = ' '
	charEscape  = '\\'
	charNewLine = '\n'
	charReturn  = '\r'
)

// Token represents a lexical unit in the ignore file.
type Token struct {
	Type    TokenType
	Literal string
}

// Pattern represents a compiled ignore rule.
type Pattern struct {
	Original string
	Path     string
	IsNegate bool
	IsDir    bool
}
