package formatting

import (
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type HumanShortFormatter struct {
}

func NewHumanShortFormatter() *HumanShortFormatter {
	return &HumanShortFormatter{}
}

func (f *HumanShortFormatter) Format(w io.Writer, v any) error {
	items, ok := v.([]fs.Item)
	if !ok {
		return fmt.Errorf("HumanShortFormatter: expected []fs.Item, got %T", v)
	}

	// Determine terminal width
	termWidth := detectTerminalWidth(w)
	if termWidth <= 0 {
		termWidth = 120 // fallback
	}

	// Determine column width based on longest display name
	colWidth := 0
	for _, it := range items {
		n := len(displayName(it))
		if n > colWidth {
			colWidth = n
		}
	}
	colWidth += 2 // padding

	cols := termWidth / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, it := range items {
		fmt.Fprintf(w, "%-*s", colWidth, displayName(it))
		if (i+1)%cols == 0 {
			fmt.Fprintln(w)
		}
	}

	fmt.Fprintln(w)
	return nil
}
