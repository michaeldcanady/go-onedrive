package ignore

// Lexer breaks the input into a stream of tokens.
type Lexer struct {
	input string
	pos   int
}

// NewLexer creates a new Lexer instance.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token in the input stream.
func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	ch := l.input[l.pos]
	switch ch {
	case charHash:
		l.pos++
		return Token{Type: TokenHash, Literal: string(ch)}
	case charBang:
		l.pos++
		return Token{Type: TokenBang, Literal: string(ch)}
	case charSlash:
		l.pos++
		return Token{Type: TokenSlash, Literal: string(ch)}
	case charSpace:
		l.pos++
		return Token{Type: TokenSpace, Literal: string(ch)}
	case '\n', '\r':
		// Handle CRLF
		if ch == '\r' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '\n' {
			l.pos++
		}
		l.pos++
		return Token{Type: TokenNewline, Literal: "\n"}
	default:
		return l.readText()
	}
}

func (l *Lexer) readText() Token {
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if isSpecial(ch) {
			break
		}
		l.pos++
	}
	return Token{Type: TokenText, Literal: l.input[start:l.pos]}
}

func isSpecial(ch byte) bool {
	switch ch {
	case charHash, charBang, charSlash, charSpace, charNewLine, charReturn:
		return true
	}
	return false
}
