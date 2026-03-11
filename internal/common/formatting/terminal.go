package formatting

import (
	"os"

	"golang.org/x/term"
)

type Terminal struct{}

func (Terminal) Width() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // fallback
	}
	return w
}
