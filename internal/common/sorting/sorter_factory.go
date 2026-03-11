package sorting

import (
	"fmt"
)

type SorterFactory struct{}

func NewSorterFactory() *SorterFactory {
	return &SorterFactory{}
}

func (f *SorterFactory) Create(opts ...SortingOption) (Sorter, error) {
	config := NewSortingOptions()

	if err := config.Apply(opts...); err != nil {
		return nil, fmt.Errorf("failed to build options: %w", err)
	}

	return NewOptionsSorter(*config), nil
}
