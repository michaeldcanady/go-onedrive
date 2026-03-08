package ignore

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"
)

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
	charHash   = '#'
	charBang   = '!'
	charSlash  = '/'
	charSpace  = ' '
	charEscape = '\\'
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

// Lexer breaks the input into a stream of tokens.
type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

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
	case charHash, charBang, charSlash, charSpace, '\n', '\r':
		return true
	}
	return false
}

// Parser converts tokens into Patterns.
type Parser struct {
	lexer *Lexer
	curr  Token
	peek  Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curr = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) Parse() []Pattern {
	var patterns []Pattern

	for p.curr.Type != TokenEOF {
		if p.curr.Type == TokenNewline || p.curr.Type == TokenSpace {
			p.nextToken()
			continue
		}

		if p.curr.Type == TokenHash {
			p.skipUntilNewline()
			continue
		}

		pattern := p.parsePattern()
		if pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns
}

func (p *Parser) skipUntilNewline() {
	for p.curr.Type != TokenEOF && p.curr.Type != TokenNewline {
		p.nextToken()
	}
}

func (p *Parser) parsePattern() *Pattern {
	pattern := &Pattern{}
	var sb strings.Builder

	if p.curr.Type == TokenBang {
		pattern.IsNegate = true
		p.nextToken()
	}

	// Read everything until newline or hash (comment)
	for p.curr.Type != TokenEOF && p.curr.Type != TokenNewline {
		if p.curr.Type == TokenHash {
			p.skipUntilNewline()
			break
		}
		sb.WriteString(p.curr.Literal)
		p.nextToken()
	}

	raw := strings.TrimSpace(sb.String())
	if raw == "" {
		return nil
	}

	pattern.Original = raw
	if strings.HasSuffix(raw, string(charSlash)) {
		pattern.IsDir = true
		pattern.Path = raw[:len(raw)-1]
	} else {
		pattern.Path = raw
	}

	return pattern
}

// Matcher evaluates paths against patterns.
type Matcher struct {
	patterns []Pattern
}

func NewMatcher(patterns []Pattern) *Matcher {
	return &Matcher{patterns: patterns}
}

func ParseReader(r io.Reader) (*Matcher, error) {
	var sb strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	l := NewLexer(sb.String())
	p := NewParser(l)
	return NewMatcher(p.Parse()), nil
}

// ShouldIgnore checks if a path matches the patterns.
func (m *Matcher) ShouldIgnore(path string, isDir bool) bool {
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimPrefix(path, "/")

	// Split path into segments to check each parent
	segments := strings.Split(path, string(charSlash))
	
	ignored := false
	currentPath := ""

	for i, segment := range segments {
		if currentPath == "" {
			currentPath = segment
		} else {
			currentPath += string(charSlash) + segment
		}

		isCurrentDir := isDir
		if i < len(segments)-1 {
			isCurrentDir = true
		}

		for _, p := range m.patterns {
			if p.IsDir && !isCurrentDir {
				continue
			}

			if matchPattern(p.Path, currentPath) {
				ignored = !p.IsNegate
			}
		}
	}

	return ignored
}

func matchPattern(pattern, path string) bool {
	// 1. Exact match
	if pattern == path {
		return true
	}

	// 2. Directory prefix match (e.g., "node_modules" matches "node_modules/foo")
	if strings.HasPrefix(path, pattern+string(charSlash)) {
		return true
	}

	// 3. Simple wildcard matching via filepath.Match
	// Note: this doesn't handle ** yet, but we can improve it.
	match, _ := filepathMatch(pattern, path)
	return match
}

// filepathMatch is a placeholder or wrapper for actual matching logic.
func filepathMatch(pattern, path string) (bool, error) {
	// For now, let's use a simple implementation.
	// We can expand this to support more complex gitignore rules later.
	parts := strings.Split(path, string(charSlash))
	for _, part := range parts {
		if ok, _ := filepath.Match(pattern, part); ok {
			return true, nil
		}
	}
	return false, nil
}
