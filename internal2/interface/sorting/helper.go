package sorting

import (
	"fmt"
	"reflect"
	"time"
)

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
			ti := vi.Interface().(time.Time)
			tj := vj.Interface().(time.Time)
			return ti.Before(tj), nil
		}
	}

	return false, fmt.Errorf("unsupported field type: %s", vi.Kind())
}
