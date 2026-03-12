package sorting

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
)

// OptionsSorter implements the Sorter interface using reflection based on the given configuration.
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
	if s.opts.Field == "" {
		return items, nil
	}

	field := s.opts.Field
	desc := s.opts.Direction == DirectionDescending

	// Validate field exists on shared.Item
	itemType := reflect.TypeOf(shared.Item{})
	if _, ok := itemType.FieldByNameFunc(func(s string) bool {
		return strings.EqualFold(s, field)
	}); !ok {
		return nil, fmt.Errorf("unknown sort field: %s", field)
	}

	sort.Slice(items, func(i, j int) bool {
		vi := reflect.ValueOf(items[i]).FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, field)
		})
		vj := reflect.ValueOf(items[j]).FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, field)
		})

		less, err := compareValues(vi, vj)
		if err != nil {
			return false
		}

		if desc {
			return !less
		}
		return less
	})

	return items, nil
}
