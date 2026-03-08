package ignore

// TokenType defines the type of a lexical token.
type TokenType int

const (
	// Special tokens
	TokenError TokenType = iota
	TokenEOF

	// Value tokens
	TokenText
	TokenComment

	// Operator/Delimiter tokens
	TokenSlash      // /
	TokenNegate     // !
	TokenStar       // *
	TokenDoubleStar // **
	TokenQuestion   // ?
	TokenOpenSet    // [
	TokenCloseSet   // ]
	TokenRange      // -
	TokenEscape     // 
	TokenNewline    // 


	TokenSpace // ' '
)
