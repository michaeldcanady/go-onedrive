package sorting

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type OptionsSorter struct {
	opts SortingOptions
}

func NewOptionsSorter(opts SortingOptions) *OptionsSorter {
	return &OptionsSorter{opts: opts}
}

func (s *OptionsSorter) Sort(items []domainfs.Item) ([]domainfs.Item, error) {
	if s.opts.Field == "" {
		return items, nil // no sorting requested
	}

	field := s.opts.Field
	desc := s.opts.Direction == DirectionDescending

	// Validate field exists on fs.Item
	itemType := reflect.TypeOf(fs.Item{})
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
