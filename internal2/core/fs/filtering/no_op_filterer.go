package filtering

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
)

// NoOpFilterer is an implementation of the Filterer interface that performs no exclusion.
type NoOpFilterer struct{}

// NewNoOpFilterer initializes a new instance of the NoOpFilterer.
func NewNoOpFilterer() *NoOpFilterer {
	return &NoOpFilterer{}
}

// Filter returns a copy of the input slice without removing any items.
func (f *NoOpFilterer) Filter(items []shared.Item) ([]shared.Item, error) {
	out := make([]shared.Item, len(items))
	copy(out, items)

	return out, nil
}
