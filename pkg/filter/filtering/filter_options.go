package filtering

import (
	"errors"
	"time"

	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
)

// FilterOptions defines the configuration criteria used by a filterer.
type FilterOptions struct {
	// ItemType restricts results to only files or only folders.
	ItemType shared.ItemType
	// IncludeAll determines whether hidden items (names starting with '.') are included.
	IncludeAll bool
	// Names restricts results to items matching any of the given glob patterns.
	Names []string
	// MinSize restricts results to items with at least the given size in bytes.
	MinSize *int64
	// MaxSize restricts results to items with at most the given size in bytes.
	MaxSize *int64
	// ModifiedBefore restricts results to items modified before the given time.
	ModifiedBefore *time.Time
	// ModifiedAfter restricts results to items modified after the given time.
	ModifiedAfter *time.Time
}

// NewFilterOptions initializes a new instance of FilterOptions with default settings.
func NewFilterOptions() *FilterOptions {
	return &FilterOptions{
		ItemType:   shared.TypeUnknown,
		IncludeAll: false,
		Names:      nil,
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
