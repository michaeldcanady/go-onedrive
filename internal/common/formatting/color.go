package formatting

import (
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"golang.org/x/term"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
	Bold   = "\033[1m"
)

// IsTerminal checks if the provided type is a terminal.
// It accepts an any type and checks if it implements an Fd() uintptr method,
// which is typical for *os.File.
func IsTerminal(w any) bool {
	if f, ok := w.(interface{ Fd() uintptr }); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

func Colorize(w io.Writer, color string, text string) string {
	if !IsTerminal(w) {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

func ColorizeItem(w io.Writer, item domain.Item) string {
	name := displayName(item)
	if item.Type == domain.ItemTypeFolder {
		return Colorize(w, Blue+Bold, name)
	}
	return name
}
