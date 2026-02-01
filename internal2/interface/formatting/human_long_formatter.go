package formatting

import (
	"fmt"
	"io"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

const (
	dateFormat = "2006-01-02 15:04"
	emptySize  = "-"
)

type HumanLongFormatter struct{}

func (f *HumanLongFormatter) Format(w io.Writer, items []domainfs.Item) error {
	for _, it := range items {
		mod := it.Modified.Format(dateFormat)

		size := emptySize
		if it.Type == domainfs.ItemTypeFile {
			size = fmt.Sprintf("%d", it.Size)
		}

		fmt.Fprintf(w, "%-20s %10s  %s\n", mod, size, displayName(it))
	}
	return nil
}
