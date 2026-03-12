package formatting

// TerminalInfo provides operations for retrieving characteristics of the user's terminal environment.
type TerminalInfo interface {
	// Width returns the current width of the terminal in characters.
	Width() int
}
