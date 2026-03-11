package formatting

import (
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type HumanShortFormatter struct {
	term TerminalInfo
}

func (f *HumanShortFormatter) Format(w io.Writer, items []domain.Item) error {
	var width = 30
	if f.term != nil {
		width = f.term.Width()
	}
	colWidth := 0

	for _, it := range items {
		if len(it.Name) > colWidth {
			colWidth = len(it.Name)
		}
	}
	colWidth += 2

	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, it := range items {
		// Calculate padding needed: colWidth - visible length
		name := displayName(it)
		padding := colWidth - len(name)
		fmt.Fprintf(w, "%s%*s", ColorizeItem(w, it), padding, "")
		if (i+1)%cols == 0 {
			fmt.Fprintln(w)
		}
	}

	fmt.Fprintln(w)
	return nil
}
