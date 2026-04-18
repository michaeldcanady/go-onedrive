package sorting

import (
	"errors"
)

// SortingOption is a functional configuration pattern for SortingOptions.
type SortingOption func(*SortingOptions) error

// WithDirection configures the last added sorting criterion with the given direction.
// If no criteria exist, it does nothing and returns an error.
func WithDirection(direction Direction) SortingOption {
	return func(o *SortingOptions) error {
		if o == nil {
			return errors.New("sorting options is nil")
		}
		if len(o.Criteria) == 0 {
			return errors.New("cannot set direction: no sorting field specified")
		}
		o.Criteria[len(o.Criteria)-1].Direction = direction
		return nil
	}
}

// WithField appends a new sorting criterion with the given field name and default ascending direction.
func WithField(name string) SortingOption {
	return func(o *SortingOptions) error {
		if o == nil {
			return errors.New("sorting options is nil")
		}
		o.Criteria = append(o.Criteria, SortingCriteria{
			Field:     name,
			Direction: DirectionAscending,
		})
		return nil
	}
}
