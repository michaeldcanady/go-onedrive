package ignore

import "strings"

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
	ch, ok := l.readChar()
	if !ok {
		return Token{Type: TokenEOF}
	}

	switch ch {
	case charHash:
		return Token{Type: TokenHash, Literal: string(ch)}
	case charBang:
		return Token{Type: TokenBang, Literal: string(ch)}
	case charSlash:
		return Token{Type: TokenSlash, Literal: string(ch)}
	case charSpace:
		return Token{Type: TokenSpace, Literal: string(ch)}
	case charEscape:
		// Handle escaped characters
		next, ok := l.readChar()
		if !ok {
			return Token{Type: TokenError, Literal: "unexpected EOF after escape"}
		}
		// If it's an escaped special char, treat it as text.
		// If it's escaped text, it's still just part of a text segment.
		// We'll peek and continue reading text if possible.
		return l.readTextFrom(string(next))
	case charNewLine, charReturn:
		// Handle CRLF
		if ch == charReturn {
			if next, ok := l.peekChar(); ok && next == charNewLine {
				l.readChar() // consume '\n'
			}
		}
		return Token{Type: TokenNewline, Literal: "\n"}
	default:
		l.unreadChar()
		return l.readText()
	}
}

// readChar reads the next character and advances the position.
func (l *Lexer) readChar() (byte, bool) {
	if l.pos >= len(l.input) {
		return 0, false
	}
	ch := l.input[l.pos]
	l.pos++
	return ch, true
}

// unreadChar moves the position back by one character.
func (l *Lexer) unreadChar() {
	if l.pos > 0 {
		l.pos--
	}
}

// peekChar allows looking at the next character without consuming it.
func (l *Lexer) peekChar() (byte, bool) {
	if l.pos >= len(l.input) {
		return 0, false
	}
	return l.input[l.pos], true
}

// readText reads characters until it encounters a special character or EOF.
func (l *Lexer) readText() Token {
	return l.readTextFrom("")
}

// readTextFrom reads characters until it encounters a special character or EOF, starting with an initial string.
func (l *Lexer) readTextFrom(initial string) Token {
	var sb strings.Builder
	sb.WriteString(initial)

	for {
		ch, ok := l.readChar()
		if !ok {
			break
		}

		if isSpecial(ch) || ch == charEscape {
			l.unreadChar()
			break
		}
		sb.WriteByte(ch)
	}

	return Token{Type: TokenText, Literal: sb.String()}
}

// isSpecial checks if a character is a special token character.
func isSpecial(ch byte) bool {
	switch ch {
	case charHash, charBang, charSlash, charSpace, charNewLine, charReturn:
		return true
	}
	return false
}
