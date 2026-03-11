package formatting

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	Name  string
	Value string
}

func TestTableFormatter_computeColumnWidths(t *testing.T) {
	tf := NewTableFormatter(
		NewColumn("Name", func(it testItem) string { return it.Name }),
		NewColumn("Val", func(it testItem) string { return it.Value }),
	)

	items := []testItem{
		{"Short", "123"},
		{"Very Long Name", "1"},
	}

	widths := tf.computeColumnWidths(items)
	assert.Equal(t, []int{14, 3}, widths) // "Very Long Name" is 14, "123" is 3 (header "Val" is also 3)
}

func TestTableFormatter_fitToTerminal(t *testing.T) {
	tf := NewTableFormatter(
		NewColumn("C1", func(it string) string { return it }),
		NewColumn("C2", func(it string) string { return it }),
	)

	widths := []int{50, 50} // Natural widths
	
	// Case 1: Fits fine
	newWidths := tf.fitToTerminal(widths, 120)
	assert.Equal(t, widths, newWidths)

	// Case 2: Needs shrinking
	// total = 50 + 50 + 2 (one separator) = 102
	// termWidth = 20
	// remaining = 20 - 2 = 18
	// minWidth = 5 each -> 10. remaining = 8.
	// Distribution: 50/100 * 8 = 4 each.
	// Final: 5 + 4 = 9 each.
	newWidths = tf.fitToTerminal(widths, 20)
	assert.Equal(t, []int{9, 9}, newWidths)
}

func TestTableFormatter_Format(t *testing.T) {
	tf := NewTableFormatter(
		NewColumn("ID", func(it testItem) string { return it.Name }),
		NewRenderColumn("Value", 
			func(it testItem) string { return it.Value },
			func(w io.Writer, it testItem) string { return "[" + it.Value + "]" },
		),
	)

	items := []testItem{
		{"A", "10"},
		{"B", "20"},
	}

	buf := new(bytes.Buffer)
	err := tf.Format(buf, items)
	assert.NoError(t, err)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Equal(t, 4, len(lines))
	assert.Contains(t, lines[0], "ID")
	assert.Contains(t, lines[0], "Value")
	assert.Contains(t, lines[2], "A")
	assert.Contains(t, lines[2], "[10]")
	assert.Contains(t, lines[3], "B")
	assert.Contains(t, lines[3], "[20]")
}
