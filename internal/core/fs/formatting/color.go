package formatting

import (
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
)

const (
	// Reset restores the terminal to its default text appearance.
	Reset = "\033[0m"
	// Red applies a red foreground color to subsequent text.
	Red = "\033[31m"
	// Green applies a green foreground color to subsequent text.
	Green = "\033[32m"
	// Yellow applies a yellow foreground color to subsequent text.
	Yellow = "\033[33m"
	// Blue applies a blue foreground color to subsequent text.
	Blue = "\033[34m"
	// Purple applies a purple foreground color to subsequent text.
	Purple = "\033[35m"
	// Cyan applies a cyan foreground color to subsequent text.
	Cyan = "\033[36m"
	// Gray applies a gray foreground color to subsequent text.
	Gray = "\033[37m"
	// White applies a white foreground color to subsequent text.
	White = "\033[97m"
	// Bold applies a bold weight to subsequent text.
	Bold = "\033[1m"
)

// Colorize wraps the given text with ANSI color codes if the writer supports terminal output.
func Colorize(w io.Writer, color string, text string) string {
	if !NewTerminal().IsTerminal(w) {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// ColorizeItem returns the display name of a filesystem item, applying color if it is a directory.
func ColorizeItem(w io.Writer, item shared.Item) string {
	name := item.Name
	if item.Type == shared.TypeFolder {
		return Colorize(w, Blue+Bold, name)
	}
	return name
}
