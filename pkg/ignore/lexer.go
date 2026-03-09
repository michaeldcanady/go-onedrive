package ignore

import (
	"unicode/utf8"
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Lexer holds the state of the scanner.
type Lexer struct {
	input  string     // the string being scanned
	start  int        // start position of this item
	pos    int        // current position in the input
	width  int        // width of last rune read from input
	line   int        // current line number
	column int        // current column number
	tokens chan Token // channel of scanned items
}

// NewLexer creates a new Lexer for the input string.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 1,
		tokens: make(chan Token),
	}
	go l.run()
	return l
}

// run lexes the input by executing state functions until the state is nil.
func (l *Lexer) run() {
	for state := lexStartOfLine; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	return <-l.tokens
}

// emit passes an item back to the client.
func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{
		Type:    t,
		Literal: l.input[l.start:l.pos],
		Position: Position{
			Line:   l.line,
			Column: l.column - (l.pos - l.start),
		},
	}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return -1 // EOF
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	l.column++
	if r == charNewline {
		l.line++
		l.column = 1
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
	l.column--
	if l.width == 1 && l.input[l.pos] == charNewline {
		l.line--
	}
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by returning nil.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{
		Type:    TokenError,
		Literal: format,
		Position: Position{
			Line:   l.line,
			Column: l.column,
		},
	}
	return nil
}

// State functions

const (
	eof = -1
)

// lexStartOfLine checks for special start-of-line tokens.
func lexStartOfLine(l *Lexer) stateFn {
	if l.pos >= len(l.input) {
		l.emit(TokenEOF)
		return nil
	}

	r := l.peek()
	switch r {
	case charHash:
		return lexComment
	case charBang:
		return lexNegate
	case charSpace:
		return lexSpace // Leading spaces are significant (ignored)
	case charNewline, charReturn:
		return lexNewline
	default:
		return lexText
	}
}

// lexText scans normal pattern text.
func lexText(l *Lexer) stateFn {
	for {
		if l.pos >= len(l.input) {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			l.emit(TokenEOF)
			return nil
		}

		r := l.peek()
		switch r {
		case charSlash:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexSlash
		case charStar:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexStar
		case charQuestion:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexQuestion
		case charOpenSet:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexOpenSet
		case charCloseSet:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexCloseSet
		case charEscape:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexEscape
		case charNewline, charReturn:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexNewline
		case charSpace:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexSpace
		default:
			l.next()
		}
	}
}

func lexComment(l *Lexer) stateFn {
	l.next() // consume #
	l.emit(TokenComment)
	// Consume until newline
	for {
		r := l.peek()
		if r == charNewline || r == charReturn || r == eof {
			break
		}
		l.next()
	}
	if l.pos > l.start {
		l.emit(TokenText)
	}
	// We do NOT consume the newline here, we let the next state handle it
	// so lexNewline is consistent.
	// But wait, if we are at newline, we should transition to lexNewline?
	// If we just return lexStartOfLine, it will see newline and go to lexNewline.
	return lexStartOfLine
}

func lexNegate(l *Lexer) stateFn {
	l.next() // consume !
	l.emit(TokenNegate)
	return lexText
}

func lexSlash(l *Lexer) stateFn {
	l.next() // consume /
	l.emit(TokenSlash)
	return lexText
}

func lexStar(l *Lexer) stateFn {
	l.next() // consume *
	if l.peek() == charStar {
		l.next()
		l.emit(TokenDoubleStar)
	} else {
		l.emit(TokenStar)
	}
	return lexText
}

func lexQuestion(l *Lexer) stateFn {
	l.next()
	l.emit(TokenQuestion)
	return lexText
}

func lexOpenSet(l *Lexer) stateFn {
	l.next()
	l.emit(TokenOpenSet)
	return lexText
}

func lexCloseSet(l *Lexer) stateFn {
	l.next()
	l.emit(TokenCloseSet)
	return lexText
}

func lexEscape(l *Lexer) stateFn {
	l.next() // consume backslash
	r := l.peek()
	if r == eof {
		l.emit(TokenText) // Trailing backslash as literal
		return lexText
	}
	l.next() // Consume escaped char
	l.emit(TokenEscape)
	return lexText
}

func lexNewline(l *Lexer) stateFn {
	r := l.next()
	if r == charReturn && l.peek() == charNewline {
		l.next()
	}
	l.emit(TokenNewline)
	return lexStartOfLine
}

func lexSpace(l *Lexer) stateFn {
	l.next()
	for l.peek() == charSpace {
		l.next()
	}
	l.emit(TokenSpace)
	return lexText
}
