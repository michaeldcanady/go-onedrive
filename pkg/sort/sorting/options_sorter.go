package sorting

import (
	"fmt"
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/order"
)

// OptionsSorter implements the Sorter interface using predefined comparators based on configuration.
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
	if len(s.opts.Criteria) == 0 {
		return items, nil
	}

	sorter := order.NewSorter(items)

	for _, c := range s.opts.Criteria {
		if c.Field == "" {
			continue
		}

		cmp, ok := s.getComparator(c.Field)
		if !ok {
			return nil, fmt.Errorf("unknown sort field: %s", c.Field)
		}

		if c.Direction == DirectionDescending {
			cmp = Reverse(cmp)
		}

		sorter.AddComparator(order.Comparator[shared.Item](cmp))
	}

	return sorter.Sort(), nil
}

// getComparator returns the appropriate ItemComparator for the given field name.
func (s *OptionsSorter) getComparator(field string) (ItemComparator, bool) {
	switch strings.ToLower(field) {
	case "name":
		return CompareByName, true
	case "size":
		return CompareBySize, true
	case "path":
		return CompareByPath, true
	case "type":
		return CompareByType, true
	case "modifiedat", "modified":
		return CompareByModifiedAt, true
	default:
		return nil, false
	}
}
