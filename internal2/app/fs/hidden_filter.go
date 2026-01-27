package fs

import (
	"strings"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type HiddenItemFilter struct {
	Enabled bool
}

func (f HiddenItemFilter) Filter(items []domainfs.Item) ([]domainfs.Item, error) {
	if !f.Enabled {
		return items, nil
	}

	out := make([]domainfs.Item, 0, len(items))
	for _, it := range items {
		if !strings.HasPrefix(it.Name, ".") {
			out = append(out, it)
		}
	}
	return out, nil
}
