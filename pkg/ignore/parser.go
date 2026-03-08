package ignore

import "strings"

// Parser parses tokens into a Document.
type Parser struct {
	tokens []Token
	pos    int
}

// NewParser creates a new Parser.
func NewParser(l *Lexer) *Parser {
	var tokens []Token
	for {
		t := l.NextToken()
		if t.Type == TokenEOF {
			break
		}
		tokens = append(tokens, t)
	}
	return &Parser{tokens: tokens}
}

// Parse parses the tokens and returns a Document.
func (p *Parser) Parse() *Document {
	doc := &Document{}
	for p.pos < len(p.tokens) {
		rule := p.parseRule()
		if rule != nil {
			doc.Rules = append(doc.Rules, *rule)
		}
	}
	return doc
}

func (p *Parser) parseRule() *Rule {
	// 1. Skip leading newlines (empty lines)
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TokenNewline {
		p.pos++
	}
	if p.pos >= len(p.tokens) {
		return nil
	}

	t := p.tokens[p.pos]
	startLine := t.Position.Line

	// 2. Check for Comment
	if t.Type == TokenComment {
		// Consume until newline
		for p.pos < len(p.tokens) && p.tokens[p.pos].Type != TokenNewline {
			p.pos++
		}
		return nil
	}

	rule := &Rule{
		Line: startLine,
	}

	// 3. Check for Negation
	if t.Type == TokenNegate {
		rule.Negate = true
		p.pos++
	}

	// 4. Skip leading spaces (Git ignores them unless escaped)
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TokenSpace {
		p.pos++
	}

	// 5. Build pattern
	var sb strings.Builder
	var lastTokenWasSlash bool

	for p.pos < len(p.tokens) {
		t = p.tokens[p.pos]
		
		if t.Type == TokenNewline {
			p.pos++
			break
		}

		shouldResetSlash := true

		switch t.Type {
		case TokenText:
			sb.WriteString(t.Literal)
		case TokenSlash:
			sb.WriteByte(charSlash)
			lastTokenWasSlash = true
			shouldResetSlash = false
		case TokenStar:
			sb.WriteByte(charStar)
		case TokenDoubleStar:
			sb.WriteString("**")
		case TokenQuestion:
			sb.WriteByte(charQuestion)
		case TokenOpenSet:
			sb.WriteByte(charOpenSet)
		case TokenCloseSet:
			sb.WriteByte(charCloseSet)
		case TokenEscape:
			if len(t.Literal) == 2 {
				char := t.Literal[1]
				if char == charStar || char == charQuestion || char == charOpenSet || char == charEscape {
					sb.WriteString(t.Literal)
				} else {
					sb.WriteByte(char)
				}
			} else {
				sb.WriteString(t.Literal)
			}
		case TokenSpace:
			// Check if trailing
			isTrailing := true
			lookahead := p.pos + 1
			for lookahead < len(p.tokens) {
				if p.tokens[lookahead].Type == TokenSpace {
					lookahead++
					continue
				}
				if p.tokens[lookahead].Type == TokenNewline || p.tokens[lookahead].Type == TokenEOF {
					break // It is trailing
				}
				isTrailing = false
				break
			}
			
			if isTrailing {
				shouldResetSlash = false // Trailing space doesn't reset "last was slash" state
			} else {
				sb.WriteString(t.Literal)
			}
		}
		
		if shouldResetSlash {
			lastTokenWasSlash = false
		}
		
		p.pos++
	}

	rule.Pattern = sb.String()
	rule.Original = rule.Pattern // In a real AST, we might reconstruct original from tokens

	if rule.Pattern == "" {
		return nil
	}

	// 6. Check for Directory Only (trailing slash)
	if lastTokenWasSlash {
		rule.DirOnly = true
		rule.Pattern = strings.TrimSuffix(rule.Pattern, string(charSlash))
	}

	return rule
}
