package formatting

import (
	"fmt"
	"io"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
)

// HumanShortFormatter implements OutputFormatter to render items in a concise, multi-column layout.
type HumanShortFormatter struct {
	// term provides information about the terminal's dimensions.
	term TerminalInfo
}

// NewHumanShortFormatter initializes a new HumanShortFormatter instance.
func NewHumanShortFormatter(term TerminalInfo) *HumanShortFormatter {
	if term == nil {
		term = NewTerminal()
	}
	return &HumanShortFormatter{term: term}
}

// Format writes the items to the output stream, automatically arranging them into columns based on terminal width.
func (f *HumanShortFormatter) Format(w io.Writer, items []any) error {
	var width = 80
	if f.term != nil {
		width = f.term.Width(w)
	}
	colWidth := 0

	for _, it := range items {
		item := it.(shared.Item)
		if len(item.Name) > colWidth {
			colWidth = len(item.Name)
		}
	}
	colWidth += 2 // Add spacing between columns

	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, it := range items {
		item := it.(shared.Item)
		coloredName := ColorizeItem(w, item)

		// Note: padding calculation here is slightly complex because coloredName contains ANSI codes
		// but we should pad based on the original name length.
		padding := colWidth - len(item.Name)
		if padding < 0 {
			padding = 0
		}

		fmt.Fprintf(w, "%s%*s", coloredName, padding, "")

		if (i+1)%cols == 0 {
			fmt.Fprintln(w)
		}
	}

	if len(items)%cols != 0 {
		fmt.Fprintln(w)
	}

	return nil
}
