package sorting

import (
	"errors"
)

// SortingCriteria defines a single field and its sort direction.
type SortingCriteria struct {
	// Field identifies the Item struct property to sort by (e.g., "Name", "Size").
	Field string
	// Direction determines whether to order items ascending or descending.
	Direction Direction
}

// SortingOptions defines the configuration criteria used by a sorter.
type SortingOptions struct {
	// Criteria is a list of sorting criteria to be applied in order.
	Criteria []SortingCriteria
}

const (
	defaultSortByField   = "Name"
	defaultSortDirection = DirectionAscending
)

// NewSortingOptions initializes a new instance of SortingOptions with default settings.
func NewSortingOptions() *SortingOptions {
	return &SortingOptions{
		Criteria: []SortingCriteria{
			{
				Field:     defaultSortByField,
				Direction: defaultSortDirection,
			},
		},
	}
}

// Apply updates the current configuration by sequentially processing the provided functional options.
func (o *SortingOptions) Apply(opts []SortingOption) error {
	if o == nil {
		return errors.New("sorting options configuration is nil")
	}

	// If options are provided, we might want to clear the default criteria
	// but only if the first option is WithField.
	// However, for backward compatibility, we'll just let Apply handle it.
	// Actually, let's clear it if ANY options are provided to avoid "Name" being always first.
	if len(opts) > 0 {
		o.Criteria = nil
	}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}
	return nil
}
