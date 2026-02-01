package sorting

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

var _ Sorter = (*PropertySorter)(nil)

type PropertySorter struct {
	opts SortingOptions
}

func NewPropertySorter(opts SortingOptions) *PropertySorter {
	return &PropertySorter{opts: opts}
}

// Sort implements [Sorter].
func (p *PropertySorter) Sort(items []fs.Item) ([]fs.Item, error) {
	field := p.opts.Field
	desc := p.opts.Direction == DirectionDescending

	// Validate field exists on fs.Item
	itemType := reflect.TypeOf(fs.Item{})
	if _, ok := itemType.FieldByName(field); !ok {
		return nil, fmt.Errorf("unknown sort field: %s", field)
	}

	sort.Slice(items, func(i, j int) bool {
		vi := reflect.ValueOf(items[i]).FieldByName(field)
		vj := reflect.ValueOf(items[j]).FieldByName(field)

		less, err := compareValues(vi, vj)
		if err != nil {
			// If comparison fails, keep original order
			return false
		}

		if desc {
			return !less
		}
		return less
	})

	return items, nil
}
