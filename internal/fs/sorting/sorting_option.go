package sorting

import (
	"errors"
)

// SortingOption is a functional configuration pattern for SortingOptions.
type SortingOption func(*SortingOptions) error

// WithDirection configures the sorter to use either ascending or descending order.
func WithDirection(direction Direction) SortingOption {
	return func(o *SortingOptions) error {
		if o == nil {
			return errors.New("sorting options is nil")
		}
		o.Direction = direction
		return nil
	}
}

// WithField configures the sorter to order items based on the given struct field name.
func WithField(name string) SortingOption {
	return func(o *SortingOptions) error {
		if o == nil {
			return errors.New("sorting options is nil")
		}
		o.Field = name
		return nil
	}
}
