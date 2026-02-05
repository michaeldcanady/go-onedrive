package filtering

import (
	"fmt"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type OptionsFilterer struct {
	opts FilterOptions
}

func NewOptionsFilterer(opts FilterOptions) *OptionsFilterer {
	return &OptionsFilterer{opts: opts}
}

func (f *OptionsFilterer) Filter(v any) error {
	items, ok := v.([]domainfs.Item)
	if !ok {
		return fmt.Errorf("expected []fs.Item, got %T", v)
	}

	// In-place compaction
	n := 0
	for _, it := range items {

		// Skip dotfiles unless IncludeAll
		if !f.opts.IncludeAll && len(it.Name) > 0 && it.Name[0] == '.' {
			continue
		}

		// Skip mismatched type
		if f.opts.ItemType != domainfs.ItemTypeUnknown &&
			it.Type != f.opts.ItemType {
			continue
		}

		// Keep item
		items[n] = it
		n++
	}

	// Reslice to new logical length
	items = items[:n]

	return nil
}
