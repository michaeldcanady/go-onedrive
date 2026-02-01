package filtering

import domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"

type NoOpFilter struct{}

func NewNoOpFilter() *NoOpFilter {
	return &NoOpFilter{}
}

func (f *NoOpFilter) Filter(items []domainfs.Item) ([]domainfs.Item, error) {
	out := make([]domainfs.Item, len(items))
	copy(out, items)

	return out, nil
}
