package list

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"golang.org/x/term"
)

var (
	driveIDColumn       = NewColumn("ID", func(item *drive.Drive) string { return item.ID })
	driveNameColumn     = NewColumn("Name", func(item *drive.Drive) string { return item.Name })
	driveTypeColumn     = NewColumn("Type", func(item *drive.Drive) string { return string(item.Type) })
	driveOwnerColumn    = NewColumn("Owner", func(item *drive.Drive) string { return string(item.Owner) })
	driveReadOnlyColumn = NewColumn("ReadOnly", func(item *drive.Drive) string { return strconv.FormatBool(item.ReadOnly) })
)

// Column defines how to extract a string value for a column.
type Column[T any] struct {
	Header string
	Value  func(item T) string
}

func NewColumn[T any](header string, value func(item T) string) Column[T] {
	return Column[T]{
		Header: header,
		Value:  value,
	}
}

type TableFormatter[T any] struct {
	Columns []Column[T]
}

func NewTableFormatter[T any](cols ...Column[T]) *TableFormatter[T] {
	return &TableFormatter[T]{Columns: cols}
}

func (tf *TableFormatter[T]) Format(w io.Writer, items any) error {
	if len(tf.Columns) == 0 {
		return fmt.Errorf("no columns defined for table formatter")
	}

	typedItems, ok := items.([]T)
	if !ok {
		return fmt.Errorf("unexpected input type: %t", items)
	}

	// Determine terminal width
	termWidth := detectTerminalWidth(w)
	if termWidth <= 0 {
		termWidth = 120 // fallback
	}

	// Compute natural column widths
	widths := tf.computeColumnWidths(typedItems)

	// Adjust widths to fit terminal
	widths = tf.fitToTerminal(widths, termWidth)

	// Write header
	tf.writeRow(w, widths, func(i int) string { return tf.Columns[i].Header })

	// Separator
	tf.writeSeparator(w, widths)

	// Rows
	for _, item := range typedItems {
		tf.writeRow(w, widths, func(i int) string {
			return tf.Columns[i].Value(item)
		})
	}

	return nil
}

func detectTerminalWidth(w io.Writer) int {
	if f, ok := w.(*os.File); ok {
		if term.IsTerminal(int(f.Fd())) {
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

func (tf *TableFormatter[T]) writeRow(w io.Writer, widths []int, get func(i int) string) {
	for i := range widths {
		text := get(i)
		text = truncate(text, widths[i])
		fmt.Fprint(w, pad(text, widths[i]))
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
