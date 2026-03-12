package filtering

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
)

// OptionsFilterer implements the Filterer interface based on a specific configuration.
type OptionsFilterer struct {
	// opts is the configuration used to include or exclude items.
	opts FilterOptions
}

// NewOptionsFilterer initializes a new OptionsFilterer instance with the provided criteria.
func NewOptionsFilterer(opts FilterOptions) *OptionsFilterer {
	return &OptionsFilterer{opts: opts}
}

// Filter evaluates a collection of items against the internal options and returns the passing set.
func (f *OptionsFilterer) Filter(items []shared.Item) ([]shared.Item, error) {
	out := make([]shared.Item, 0, len(items))

	for _, it := range items {
		// Filter out hidden items if not included.
		if !f.opts.IncludeAll && len(it.Name) > 0 && it.Name[0] == '.' {
			continue
		}

		// Filter by specific item type if requested.
		if f.opts.ItemType != shared.TypeUnknown {
			if it.Type != f.opts.ItemType {
				continue
			}
		}

		out = append(out, it)
	}

	return out, nil
}
