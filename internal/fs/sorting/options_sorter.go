package sorting

import (
	"sort"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
)

// OptionsSorter implements the Sorter interface using composable comparators.
type OptionsSorter struct {
	// opts is the configuration used to determine the field and direction for sorting.
	opts SortingOptions
}

// NewOptionsSorter initializes a new OptionsSorter instance with the provided criteria.
func NewOptionsSorter(opts SortingOptions) *OptionsSorter {
	return &OptionsSorter{opts: opts}
}

// Sort orders a collection of items according to the internal options and returns the resulting slice.
func (s *OptionsSorter) Sort(items []shared.Item) ([]shared.Item, error) {
	comparator := NewItemComparator(s.opts)

	sort.Slice(items, func(i, j int) bool {
		return comparator.Less(items[i], items[j])
	})

	return items, nil
}
