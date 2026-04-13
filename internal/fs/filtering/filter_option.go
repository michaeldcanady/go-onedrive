package filtering

import (
	"errors"
	"time"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
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

// WithName restricts results to items matching any of the given glob patterns.
func WithName(pattern string) FilterOption {
	return func(config *FilterOptions) error {
		config.Names = append(config.Names, pattern)
		return nil
	}
}

// WithMinSize restricts results to items with at least the given size in bytes.
func WithMinSize(size int64) FilterOption {
	return func(config *FilterOptions) error {
		config.MinSize = &size
		return nil
	}
}

// WithMaxSize restricts results to items with at most the given size in bytes.
func WithMaxSize(size int64) FilterOption {
	return func(config *FilterOptions) error {
		config.MaxSize = &size
		return nil
	}
}

// WithModifiedBefore restricts results to items modified before the given time.
func WithModifiedBefore(t time.Time) FilterOption {
	return func(config *FilterOptions) error {
		config.ModifiedBefore = &t
		return nil
	}
}

// WithModifiedAfter restricts results to items modified after the given time.
func WithModifiedAfter(t time.Time) FilterOption {
	return func(config *FilterOptions) error {
		config.ModifiedAfter = &t
		return nil
	}
}
