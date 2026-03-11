package formatting

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Column defines how to extract a string value for a column.
type Column[T any] struct {
	Header string
	Value  func(item T) string
	Render func(w io.Writer, item T) string
}

func NewColumn[T any](header string, value func(item T) string) Column[T] {
	return Column[T]{
		Header: header,
		Value:  value,
	}
}

func NewRenderColumn[T any](header string, value func(item T) string, render func(w io.Writer, item T) string) Column[T] {
	return Column[T]{
		Header: header,
		Value:  value,
		Render: render,
	}
}

type TableFormatter[T any] struct {
	Columns  []Column[T]
	Truncate bool
}

func NewTableFormatter[T any](cols ...Column[T]) *TableFormatter[T] {
	return &TableFormatter[T]{Columns: cols}
}

func (tf *TableFormatter[T]) WithTruncate(truncate bool) *TableFormatter[T] {
	tf.Truncate = truncate
	return tf
}

func (tf *TableFormatter[T]) Format(w io.Writer, items []T) error {
	if len(tf.Columns) == 0 {
		return fmt.Errorf("no columns defined for table formatter")
	}

	// Determine terminal width
	termWidth := detectTerminalWidth(w)
	if termWidth <= 0 {
		termWidth = 120 // fallback
	}

	// Compute natural column widths
	widths := tf.computeColumnWidths(items)

	// Adjust widths to fit terminal
	widths = tf.fitToTerminal(widths, termWidth)

	// Write header
	tf.writeRow(w, widths, func(i int) (string, int) {
		h := tf.Columns[i].Header
		if tf.Truncate && len(h) > widths[i] {
			h = truncate(h, widths[i])
		}
		return h, len(h)
	})

	// Separator
	tf.writeSeparator(w, widths)

	// Rows
	for _, item := range items {
		tf.writeRow(w, widths, func(i int) (string, int) {
			val := tf.Columns[i].Value(item)
			if tf.Truncate && len(val) > widths[i] {
				val = truncate(val, widths[i])
			}
			visibleLen := len(val)
			if tf.Columns[i].Render != nil {
				// We pass a modified writer or use a different approach if Render needs the truncated value.
				// Currently Render takes the full item T.
				// If we want Render to respect truncation, we might need to change its signature
				// or assume it uses the Value function internally.
				return tf.Columns[i].Render(w, item), visibleLen
			}
			return val, visibleLen
		})
	}

	return nil
}

func detectTerminalWidth(w io.Writer) int {
	if IsTerminal(w) {
		if f, ok := w.(*os.File); ok {
			width, _, err := term.GetSize(int(f.Fd()))
			if err == nil {
				return width
			}
		}
	}
	return -1
}

func (tf *TableFormatter[T]) computeColumnWidths(items []T) []int {
	widths := make([]int, len(tf.Columns))

	// Start with header widths
	for i, col := range tf.Columns {
		widths[i] = len(col.Header)
	}

	// Expand based on values
	for _, item := range items {
		for i, col := range tf.Columns {
			val := col.Value(item)
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
	}

	return widths
}

func (tf *TableFormatter[T]) fitToTerminal(widths []int, termWidth int) []int {
	// Compute total width including spacing
	total := 0
	for _, w := range widths {
		total += w
	}
	total += (len(widths) - 1) * 2 // spaces between columns

	if total <= termWidth {
		return widths // fits fine
	}

	// Need to shrink columns
	remaining := termWidth - (len(widths)-1)*2
	newWidths := make([]int, len(widths))

	// Strategy:
	// 1. Give each column at least 5 chars
	// 2. Distribute remaining width proportionally

	minWidth := 5
	for i := range widths {
		newWidths[i] = minWidth
		remaining -= minWidth
	}

	if remaining < 0 {
		return newWidths // extremely narrow terminal
	}

	// Distribute remaining width proportionally
	totalNatural := 0
	for _, w := range widths {
		totalNatural += w
	}

	for i := range widths {
		extra := (widths[i] * remaining) / totalNatural
		newWidths[i] += extra
	}

	return newWidths
}

func (tf *TableFormatter[T]) writeRow(w io.Writer, widths []int, get func(i int) (string, int)) {
	for i := range widths {
		text, visibleLen := get(i)
		// Truncate based on visible length if necessary (simplistic)
		// For now, assume column width is sufficient or handle truncation carefully
		// truncate() only handles string length, not visible length with colors.
		// If we support colors, truncation is hard without parsing ANSI.
		// Let's assume we don't truncate colored output for now or just truncate plain text if render is not used.

		// Padding
		padding := widths[i] - visibleLen
		if padding < 0 {
			padding = 0
		}

		fmt.Fprint(w, text)
		fmt.Fprint(w, strings.Repeat(" ", padding))

		if i < len(widths)-1 {
			fmt.Fprint(w, "  ")
		}
	}
	fmt.Fprintln(w)
}

func (tf *TableFormatter[T]) writeSeparator(w io.Writer, widths []int) {
	for i, width := range widths {
		fmt.Fprint(w, strings.Repeat("-", width))
		if i < len(widths)-1 {
			fmt.Fprint(w, "  ")
		}
	}
	fmt.Fprintln(w)
}

func truncate(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width <= 1 {
		return s[:width]
	}
	return s[:width-1] + "…"
}

func pad(s string, width int) string {
	if len(s) < width {
		return s + strings.Repeat(" ", width-len(s))
	}
	return s
}
