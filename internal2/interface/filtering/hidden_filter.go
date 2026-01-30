package filtering

import (
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type HiddenFilterer struct {
	IncludeHidden bool
}

func NewHiddenFilterer() *HiddenFilterer {
	return &HiddenFilterer{}
}

func (f *HiddenFilterer) Filter(items []domainfs.Item) ([]domainfs.Item, error) {
	out := make([]domainfs.Item, 0, len(items))
	for _, it := range items {
		if len(it.Name) > 0 && it.Name[0] == '.' {
			continue
		}
		out = append(out, it)
	}
	return out, nil
}
