package filtering

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
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
	specification := NewFilterSpec(f.opts)

	for _, item := range items {
		if specification.IsSatisfiedBy(item) {
			out = append(out, item)
		}
	}

	return out, nil
}
