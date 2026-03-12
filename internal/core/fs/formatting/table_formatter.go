package formatting

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Column defines the presentation and data extraction logic for a single table field.
type Column struct {
	// Header is the title displayed at the top of the column.
	Header string
	// Value returns the raw string representation of the field for the given item.
	Value func(item any) string
	// Render optionally provides a colored or stylized representation of the field.
	Render func(w io.Writer, item any) string
}

// NewColumn initializes a basic data column with a header and a value extraction function.
func NewColumn(header string, value func(item any) string) Column {
	return Column{
		Header: header,
		Value:  value,
	}
}

// NewRenderColumn initializes a column that supports custom rendering (e.g., ANSI colors).
func NewRenderColumn(header string, value func(item any) string, render func(w io.Writer, item any) string) Column {
	return Column{
		Header: header,
		Value:  value,
		Render: render,
	}
}

// TableFormatter renders a collection of items as an aligned text table.
type TableFormatter struct {
	// Columns defines the fields to include in the table.
	Columns []Column
	// Truncate determines whether long values should be shortened to fit the terminal width.
	Truncate bool
}

// NewTableFormatter initializes a new TableFormatter with the specified columns.
func NewTableFormatter(cols ...Column) *TableFormatter {
	return &TableFormatter{Columns: cols}
}

// WithTruncate configures the formatter to shorten overflowing column values.
func (tf *TableFormatter) WithTruncate(truncate bool) *TableFormatter {
	tf.Truncate = truncate
	return tf
}

// Format writes the items as an ASCII table to the provided writer, adjusting for terminal width.
func (tf *TableFormatter) Format(w io.Writer, items []any) error {
	if len(tf.Columns) == 0 {
		return fmt.Errorf("no columns defined for table formatter")
	}

	termWidth := tf.detectTerminalWidth(w)
	if termWidth <= 0 {
		termWidth = 120
	}

	widths := tf.computeColumnWidths(items)
	widths = tf.fitToTerminal(widths, termWidth)

	tf.writeRow(w, widths, func(i int) (string, int) {
		h := tf.Columns[i].Header
		if tf.Truncate && len(h) > widths[i] {
			h = tf.truncate(h, widths[i])
		}
		return h, len(h)
	})

	tf.writeSeparator(w, widths)

	for _, item := range items {
		tf.writeRow(w, widths, func(i int) (string, int) {
			val := tf.Columns[i].Value(item)
			if tf.Truncate && len(val) > widths[i] {
				val = tf.truncate(val, widths[i])
			}
			visibleLen := len(val)
			if tf.Columns[i].Render != nil {
				return tf.Columns[i].Render(w, item), visibleLen
			}
			return val, visibleLen
		})
	}

	return nil
}

func (tf *TableFormatter) detectTerminalWidth(w io.Writer) int {
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

func (tf *TableFormatter) computeColumnWidths(items []any) []int {
	widths := make([]int, len(tf.Columns))
	for i, col := range tf.Columns {
		widths[i] = len(col.Header)
	}
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

func (tf *TableFormatter) fitToTerminal(widths []int, termWidth int) []int {
	total := 0
	for _, w := range widths {
		total += w
	}
	total += (len(widths) - 1) * 2

	if total <= termWidth {
		return widths
	}

	remaining := termWidth - (len(widths)-1)*2
	newWidths := make([]int, len(widths))
	minWidth := 5
	for i := range widths {
		newWidths[i] = minWidth
		remaining -= minWidth
	}

	if remaining < 0 {
		return newWidths
	}

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

func (tf *TableFormatter) writeRow(w io.Writer, widths []int, get func(i int) (string, int)) {
	for i := range widths {
		text, visibleLen := get(i)
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

func (tf *TableFormatter) writeSeparator(w io.Writer, widths []int) {
	for i, width := range widths {
		fmt.Fprint(w, strings.Repeat("-", width))
		if i < len(widths)-1 {
			fmt.Fprint(w, "  ")
		}
	}
	fmt.Fprintln(w)
}

func (tf *TableFormatter) truncate(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width <= 1 {
		return s[:width]
	}
	return s[:width-1] + "…"
}
