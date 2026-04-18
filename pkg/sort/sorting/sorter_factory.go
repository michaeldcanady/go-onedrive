package sorting

import (
	"fmt"
)

// SorterFactory provides operations for initializing configured sorter instances.
type SorterFactory struct{}

// NewSorterFactory initializes a new instance of the SorterFactory.
func NewSorterFactory() *SorterFactory {
	return &SorterFactory{}
}

// Create initializes and configures a new sorter with the provided functional options.
func (f *SorterFactory) Create(opts ...SortingOption) (Sorter, error) {
	config := NewSortingOptions()

	if err := config.Apply(opts); err != nil {
		return nil, fmt.Errorf("failed to build sorting options: %w", err)
	}

	return NewOptionsSorter(*config), nil
}
