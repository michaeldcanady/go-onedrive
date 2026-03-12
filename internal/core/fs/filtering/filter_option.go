package filtering

import (
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
)

// FilterOption is a functional configuration pattern for FilterOptions.
type FilterOption func(*FilterOptions) error

// WithItemType restricts the filter results to only items matching the given type.
func WithItemType(itemType shared.ItemType) FilterOption {
	if itemType == shared.TypeUnknown {
		return func(_ *FilterOptions) error {
			return errors.New("filtered item type is unknown")
		}
	}
	return func(config *FilterOptions) error {
		config.ItemType = itemType
		return nil
	}
}

// IncludeAll configures the filter to include hidden items (names starting with '.').
func IncludeAll() FilterOption {
	return func(config *FilterOptions) error {
		config.IncludeAll = true
		return nil
	}
}

// ExcludeHidden configures the filter to omit hidden items (names starting with '.').
func ExcludeHidden() FilterOption {
	return func(config *FilterOptions) error {
		config.IncludeAll = false
		return nil
	}
}
