package formatting

import (
	"os"

	"golang.org/x/term"
)

// Terminal provides a concrete implementation of TerminalInfo using the standard operating system terminal.
type Terminal struct{}

// Width retrieves the horizontal character count of the standard output terminal.
// It returns a fallback of 80 if the width cannot be determined.
func (Terminal) Width() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return w
}
