package formatting

import (
	"io"
	"os"

	"golang.org/x/term"
)

// TerminalInfo provides the operations for querying terminal state and dimensions.
type TerminalInfo interface {
	// IsTerminal checks if the provided writer is connected to an interactive terminal session.
	IsTerminal(w io.Writer) bool
	// Width retrieves the horizontal character count of the terminal associated with the writer.
	Width(w io.Writer) int
}

// Terminal provides a concrete implementation of TerminalInfo using standard operating system interfaces.
type Terminal struct{}

// NewTerminal initializes a new instance of the Terminal.
func NewTerminal() Terminal {
	return Terminal{}
}

// IsTerminal determines if the provided writer is connected to an interactive terminal session.
func (Terminal) IsTerminal(w io.Writer) bool {
	if f, ok := w.(interface{ Fd() uintptr }); ok {
		fd := f.Fd()
		if fd > uintptr(int(^uint(0)>>1)) {
			return false
		}
		return term.IsTerminal(int(fd))
	}
	return false
}

// Width retrieves the horizontal character count of the terminal associated with the writer.
// It returns a fallback of 80 if the width cannot be determined.
func (t Terminal) Width(w io.Writer) int {
	if !t.IsTerminal(w) {
		return 80
	}

	if f, ok := w.(interface{ Fd() uintptr }); ok {
		fd := f.Fd()
		if fd <= uintptr(int(^uint(0)>>1)) {
			width, _, err := term.GetSize(int(fd))
			if err == nil {
				return width
			}
		}
	}

	// Fallback to standard output if the writer's width can't be resolved
	fd := os.Stdout.Fd()
	if fd <= uintptr(int(^uint(0)>>1)) {
		if width, _, err := term.GetSize(int(fd)); err == nil {
			return width
		}
	}

	return 80
}
