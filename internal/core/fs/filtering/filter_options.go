package filtering

import (
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
)

// FilterOptions defines the configuration criteria used by a filterer.
type FilterOptions struct {
	// ItemType restricts results to only files or only folders.
	ItemType shared.ItemType
	// IncludeAll determines whether hidden items (names starting with '.') are included.
	IncludeAll bool
}

// NewFilterOptions initializes a new instance of FilterOptions with default settings.
func NewFilterOptions() *FilterOptions {
	return &FilterOptions{
		ItemType: shared.TypeUnknown,
	}
}

// Apply updates the current configuration by sequentially processing the provided functional options.
func (o *FilterOptions) Apply(opts []FilterOption) error {
	if o == nil {
		return errors.New("filter options configuration is nil")
	}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}
	return nil
}
