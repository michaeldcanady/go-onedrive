package formatting

import (
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
)

// HumanShortFormatter implements OutputFormatter to render items in a concise, multi-column layout.
type HumanShortFormatter struct {
	// term provides information about the terminal's dimensions.
	term TerminalInfo
}

// NewHumanShortFormatter initializes a new HumanShortFormatter instance.
func NewHumanShortFormatter(term TerminalInfo) *HumanShortFormatter {
	return &HumanShortFormatter{term: term}
}

// Format writes the items to the output stream, automatically arranging them into columns based on terminal width.
func (f *HumanShortFormatter) Format(w io.Writer, items []any) error {
	var width = 30
	if f.term != nil {
		width = f.term.Width()
	}
	colWidth := 0

	for _, it := range items {
		item := it.(shared.Item)
		if len(item.Name) > colWidth {
			colWidth = len(item.Name)
		}
	}
	colWidth += 2

	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, it := range items {
		item := it.(shared.Item)
		name := item.Name
		padding := colWidth - len(name)
		fmt.Fprintf(w, "%s%*s", ColorizeItem(w, item), padding, "")
		if (i+1)%cols == 0 {
			fmt.Fprintln(w)
		}
	}

	if len(items)%cols != 0 {
		fmt.Fprintln(w)
	}

	return nil
}
