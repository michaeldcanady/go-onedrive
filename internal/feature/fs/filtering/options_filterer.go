package filtering

import (
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/feature/fs"
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

// Filter evaluates each item against the configured options and returns only the matches.
func (f *OptionsFilterer) Filter(items []shared.Item) ([]shared.Item, error) {
	var out []shared.Item

	for _, item := range items {
		// 1. Check for hidden items
		if !f.opts.IncludeAll && strings.HasPrefix(item.Name, ".") {
			continue
		}

		// 2. Check for item type
		if f.opts.ItemType != shared.TypeUnknown && item.Type != f.opts.ItemType {
			continue
		}

		out = append(out, item)
	}

	return out, nil
}
