package filtering

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

// OptionsFilterer implements the Filterer interface based on a specific configuration.
type OptionsFilterer struct {
	// opts is the configuration used to include or exclude items.
	opts spec.Specification[shared.Item]
}

// NewOptionsFilterer initializes a new OptionsFilterer instance with the provided criteria.
func NewOptionsFilterer(opts FilterOptions) *OptionsFilterer {
	var specs []spec.Specification[shared.Item]

	// 1. Hidden Filter
	if !opts.IncludeAll {
		// If IncludeAll is false, we want to EXCLUDE hidden items.
		// HiddenFilter{hidden: false} returns true for items that are NOT hidden.
		specs = append(specs, &HiddenFilter{hidden: false})
	}

	// 2. Item Type Filter
	if opts.ItemType != shared.TypeUnknown {
		specs = append(specs, &ItemTypeFilter{itemType: opts.ItemType})
	}

	// 3. Name Filter
	if len(opts.Names) > 0 {
		specs = append(specs, NewNameFilter(opts.Names))
	}

	// 4. Size Filter
	if opts.MinSize != nil || opts.MaxSize != nil {
		specs = append(specs, NewSizeFilter(opts.MinSize, opts.MaxSize))
	}

	// 5. Date Filter
	if opts.ModifiedBefore != nil || opts.ModifiedAfter != nil {
		specs = append(specs, NewDateFilter(opts.ModifiedBefore, opts.ModifiedAfter))
	}

	filter := spec.AndAll(specs...)

	return &OptionsFilterer{opts: filter}
}

// Filter evaluates each item against the configured options and returns only the matches.
func (f *OptionsFilterer) Filter(items []shared.Item) ([]shared.Item, error) {
	var out []shared.Item

	for _, item := range items {
		if f.opts.IsSatisfiedBy(item) {
			out = append(out, item)
		}
	}

	return out, nil
}
