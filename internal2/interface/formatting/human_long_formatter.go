package formatting

import (
	"io"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

const (
	dateFormat = "2006-01-02 15:04"
	emptySize  = "-"
)

type HumanLongFormatter struct{}

func (f *HumanLongFormatter) Format(w io.Writer, items []domainfs.Item) error {
	formatter := NewTableFormatter(
		NewColumn("Modified", func(it domainfs.Item) string { return it.Modified.Format(dateFormat) }),
		NewColumn("Size", func(it domainfs.Item) string {
			if it.Type == domainfs.ItemTypeFile {
				return FormatSize(it.Size)
			}
			return emptySize
		}),
		NewRenderColumn("Name",
			func(it domainfs.Item) string { return displayName(it) },
			func(w io.Writer, it domainfs.Item) string { return ColorizeItem(w, it) },
		),
	)

	return formatter.Format(w, items)
}
