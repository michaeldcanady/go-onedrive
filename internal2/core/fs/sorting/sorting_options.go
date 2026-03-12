package sorting

import (
	"errors"
)

// SortingOptions defines the configuration criteria used by a sorter.
type SortingOptions struct {
	// Direction determines whether to order items ascending or descending.
	Direction Direction
	// Field identifies the Item struct property to sort by (e.g., "Name", "Size").
	Field string
}

// NewSortingOptions initializes a new instance of SortingOptions with default settings.
func NewSortingOptions() *SortingOptions {
	return &SortingOptions{
		Direction: DirectionAscending,
		Field:     "Name",
	}
}

// Apply updates the current configuration by sequentially processing the provided functional options.
func (o *SortingOptions) Apply(opts []SortingOption) error {
	if o == nil {
		return errors.New("sorting options configuration is nil")
	}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}
	return nil
}
