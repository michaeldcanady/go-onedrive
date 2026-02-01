package sorting

import (
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

var _ Sorter = (*NoOpSorter)(nil)

type NoOpSorter struct{}

func NewNoOpSorter() *NoOpSorter {
	return &NoOpSorter{}
}

func (s *NoOpSorter) Sort(items []domainfs.Item) ([]domainfs.Item, error) {
	out := make([]domainfs.Item, len(items))
	copy(out, items)
	return out, nil
}
