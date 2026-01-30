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

// compareValues compares two reflect.Values and returns true if vi < vj.
func compareValues(vi, vj reflect.Value) (bool, error) {
	switch vi.Kind() {
	case reflect.String:
		return vi.String() < vj.String(), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vi.Int() < vj.Int(), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vi.Uint() < vj.Uint(), nil

	case reflect.Float32, reflect.Float64:
		return vi.Float() < vj.Float(), nil

	case reflect.Bool:
		// false < true
		return !vi.Bool() && vj.Bool(), nil

	case reflect.Struct:
		// TODO: find way to support custom types?
		// Special case: time.Time
		if vi.Type().String() == "time.Time" {
			ti := vi.Interface().(interface{ Before(t interface{}) bool })
			tj := vj.Interface()
			return ti.Before(tj), nil
		}
	}

	return false, fmt.Errorf("unsupported field type: %s", vi.Kind())
}
