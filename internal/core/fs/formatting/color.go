package formatting

import (
	"io"

	"github.com/fatih/color"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
)

var (
	// Blue color for directories
	blue = color.New(color.FgBlue, color.Bold)
	// Default color for files
	defaultColor = color.New()
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
