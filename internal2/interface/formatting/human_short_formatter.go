package formatting

import (
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type HumanShortFormatter struct {
	term TerminalInfo
}

func (f *HumanShortFormatter) Format(w io.Writer, items []fs.Item) error {
	var width = 30
	if f.term != nil {
		width = f.term.Width()
	}
	colWidth := 0

	for _, it := range items {
		if len(displayName(it)) > colWidth {
			colWidth = len(displayName(it))
		}
	}
	colWidth += 2

	cols := width / colWidth
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
