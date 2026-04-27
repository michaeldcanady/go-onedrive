package formatting

import (
	"io"

	"github.com/fatih/color"
	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
)

var (
	// Blue color for directories
	blue = color.New(color.FgBlue, color.Bold)
)

// Colorize wraps the given text with ANSI color codes if the writer supports terminal output.
func Colorize(w io.Writer, c *color.Color, text string) string {
	if color.NoColor || !NewTerminal().IsTerminal(w) {
		return text
	}
	return c.Sprint(text)
}

// ColorizeItem returns the display name of a filesystem item, applying color if it is a directory.
func ColorizeItem(w io.Writer, item shared.Item) string {
	name := item.Name
	if item.Type == shared.TypeFolder {
		return Colorize(w, blue, name)
	}
	return name
}
