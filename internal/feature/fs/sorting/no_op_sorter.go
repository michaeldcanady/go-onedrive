package sorting

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/feature/fs"
)

// NoOpSorter is an implementation of the Sorter interface that performs no ordering.
type NoOpSorter struct{}

// NewNoOpSorter initializes a new instance of the NoOpSorter.
func NewNoOpSorter() *NoOpSorter {
	return &NoOpSorter{}
}

// Sort returns a copy of the input slice without changing the order of items.
func (s *NoOpSorter) Sort(items []shared.Item) ([]shared.Item, error) {
	out := make([]shared.Item, len(items))
	copy(out, items)
	return out, nil
}
