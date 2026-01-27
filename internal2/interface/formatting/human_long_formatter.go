package formatting

import (
	"fmt"
	"io"
	"sort"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type HumanLongFormatter struct{}

func (f *HumanLongFormatter) Format(w io.Writer, items []domainfs.Item) error {
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })

	for _, it := range items {
		mod := it.Modified.Format("2006-01-02 15:04")

		size := "-"
		if it.Type == domainfs.ItemTypeFile {
			size = fmt.Sprintf("%d", it.Size)
		}

		fmt.Fprintf(w, "%-20s %10s  %s\n", mod, size, displayName(it))
	}
	return nil
}
