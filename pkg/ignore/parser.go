package ignore

import (
	"strings"
)

// Parser converts tokens into Patterns.
type Parser struct {
	lexer *Lexer
	curr  Token
	peek  Token
}

// NewParser creates a new Parser instance.
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

// Parse converts the token stream into a slice of Patterns.
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
	// Trim leading slash to make it root-relative in our simplified matcher
	path := strings.TrimPrefix(raw, string(charSlash))

	if strings.HasSuffix(path, string(charSlash)) {
		pattern.IsDir = true
		pattern.Path = path[:len(path)-1]
	} else {
		pattern.Path = path
	}

	return pattern
}
