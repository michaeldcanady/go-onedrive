package ignore

import "fmt"

// Token represents a lexical unit.
type Token struct {
	Type     TokenType
	Literal  string
	Position Position
}

func (t Token) String() string {
	switch t.Type {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return t.Literal
	}
	if len(t.Literal) > 10 {
		return fmt.Sprintf("%.10q...", t.Literal)
	}
	return fmt.Sprintf("%q", t.Literal)
}
