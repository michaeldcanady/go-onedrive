package sorting

import (
	"reflect"
	"sort"

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

	sort.Slice(items, func(i, j int) bool {
		vi := reflect.ValueOf(items[i]).FieldByName(field)
		vj := reflect.ValueOf(items[j]).FieldByName(field)

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
