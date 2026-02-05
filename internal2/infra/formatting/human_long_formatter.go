package formatting

import (
	"fmt"
	"io"
	"reflect"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

const (
	dateFormat = "2006-01-02 15:04"
	emptySize  = "-"
)

type HumanLongFormatter struct{}

func (f *HumanLongFormatter) Format(w io.Writer, v any) error {
	// Accept both []domainfs.Item and []*domainfs.Item
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("expects a slice, got %T", v)
	}

	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i).Interface()

		var it domainfs.Item

		switch typed := elem.(type) {
		case domainfs.Item:
			it = typed
		case *domainfs.Item:
			it = *typed
		default:
			return fmt.Errorf("unsupported element type %T", elem)
		}

		mod := it.Modified.Format(dateFormat)

		size := emptySize
		if it.Type == domainfs.ItemTypeFile {
			size = fmt.Sprintf("%d", it.Size)
		}

		fmt.Fprintf(w, "%-20s %10s  %s\n", mod, size, displayName(it))
	}

	return nil
}
