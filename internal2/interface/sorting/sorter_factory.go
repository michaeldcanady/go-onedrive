package sorting

import (
	"errors"
	"fmt"
	"strings"
)

type SorterFactory struct{}

func NewSorterFactory() *SorterFactory {
	return &SorterFactory{}
}

func (f *SorterFactory) Create(sortType string, opts ...SortingOption) (Sorter, error) {
	config := NewSortingOptions()

	if err := config.Apply(opts...); err != nil {
		return nil, fmt.Errorf("failed to build options: %w", err)
	}

	switch strings.ToLower(sortType) {
	case "property":
		return NewPropertySorter(*config), nil
	}
	return nil, errors.New("unknown sort type")
}
