package formatting

import (
	"io"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
)

const (
	dateFormat = "2006-01-02 15:04"
	emptySize  = "-"
)

// HumanLongFormatter implements OutputFormatter to render items in a detailed table with metadata.
type HumanLongFormatter struct{}

// NewHumanLongFormatter initializes a new HumanLongFormatter instance.
func NewHumanLongFormatter() *HumanLongFormatter {
	return &HumanLongFormatter{}
}

// Format writes the items to the output stream as a table containing modification time, size, and name.
func (f *HumanLongFormatter) Format(w io.Writer, items []any) error {
	formatter := NewTableFormatter(
		NewColumn("Modified", func(it any) string {
			item := it.(shared.Item)
			return item.ModifiedAt.Format(dateFormat)
		}),
		NewColumn("Size", func(it any) string {
			item := it.(shared.Item)
			if item.Type == shared.TypeFile {
				return FormatSize(item.Size)
			}
			return emptySize
		}),
		NewRenderColumn("Name",
			func(it any) string { return it.(shared.Item).Name },
			func(w io.Writer, it any) string { return ColorizeItem(w, it.(shared.Item)) },
		),
	)

	return formatter.Format(w, items)
}
