package filtering

import domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"

type OptionsFilterer struct {
	opts FilterOptions
}

func NewOptionsFilterer(opts FilterOptions) *OptionsFilterer {
	return &OptionsFilterer{opts: opts}
}

func (f *OptionsFilterer) Filter(items []domainfs.Item) ([]domainfs.Item, error) {
	out := make([]domainfs.Item, 0, len(items))

	for _, it := range items {

		if !f.opts.IncludeAll && len(it.Name) > 0 && it.Name[0] == '.' {
			continue
		}

		if f.opts.ItemType != domainfs.ItemTypeUnknown {
			if it.Type != f.opts.ItemType {
				continue
			}
		}

		out = append(out, it)
	}

	return out, nil
}
