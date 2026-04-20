package formatting

import (
	"fmt"
	"io"
	"reflect"
	"strings"
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
	// terminal provides terminal-related information and detection.
	terminal TerminalInfo
}

// NewTableFormatter initializes a new TableFormatter with the specified columns.
func NewTableFormatter(cols ...Column) *TableFormatter {
	return &TableFormatter{
		Columns:  cols,
		terminal: NewTerminal(), // Default to NewTerminal()
	}
}

// WithTerminal configures the formatter to use the provided terminal information.
func (tf *TableFormatter) WithTerminal(term TerminalInfo) *TableFormatter {
	if term != nil {
		tf.terminal = term
	}
	return tf
}

// WithTruncate configures the formatter to shorten overflowing column values.
func (tf *TableFormatter) WithTruncate(truncate bool) *TableFormatter {
	tf.Truncate = truncate
	return tf
}

// Format writes the items as an ASCII table to the provided writer, adjusting for terminal width.
func (tf *TableFormatter) Format(w io.Writer, items []any) error {
	cols := tf.Columns
	if len(cols) == 0 && len(items) > 0 {
		t := reflect.TypeOf(items[0])
		if registered, ok := GlobalRegistry.GetTable(t); ok {
			cols = registered
		}
	}

	if len(cols) == 0 {
		return fmt.Errorf("no columns defined for table formatter")
	}

	// Use temporary columns for calculation and rendering if dynamically resolved
	currentCols := cols

	termWidth := tf.terminal.Width(w)

	widths := tf.computeColumnWidths(items, currentCols)
	widths = tf.fitToTerminal(widths, termWidth)

	tf.writeRow(w, widths, func(i int) (string, int) {
		h := currentCols[i].Header
		if tf.Truncate && len(h) > widths[i] {
			h = tf.truncate(h, widths[i])
		}
		return h, len(h)
	})

	tf.writeSeparator(w, widths)

	for _, item := range items {
		tf.writeRow(w, widths, func(i int) (string, int) {
			val := currentCols[i].Value(item)
			if tf.Truncate && len(val) > widths[i] {
				val = tf.truncate(val, widths[i])
			}
			visibleLen := len(val)
			if currentCols[i].Render != nil {
				return currentCols[i].Render(w, item), visibleLen
			}
			return val, visibleLen
		})
	}

	return nil
}

func (tf *TableFormatter) computeColumnWidths(items []any, cols []Column) []int {
	widths := make([]int, len(cols))
	for i, col := range cols {
		widths[i] = len(col.Header)
	}
	for _, item := range items {
		for i, col := range cols {
			val := col.Value(item)
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
	}
	return widths
}

func (tf *TableFormatter) fitToTerminal(widths []int, termWidth int) []int {
	totalNatural := 0
	for _, w := range widths {
		totalNatural += w
	}
	// Add padding between columns (2 spaces per separator)
	total := totalNatural + (len(widths)-1)*2

	if total <= termWidth {
		return widths
	}

	// Calculate how much space we actually have for the columns themselves
	available := termWidth - (len(widths)-1)*2
	if available < len(widths)*minColWidth {
		available = len(widths) * minColWidth
	}

	newWidths := make([]int, len(widths))
	// Distribute available space proportionally
	for i, w := range widths {
		newWidths[i] = (w * available) / totalNatural
		if newWidths[i] < minColWidth {
			newWidths[i] = minColWidth
		}
	}

	return newWidths
}

const minColWidth = 5

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
