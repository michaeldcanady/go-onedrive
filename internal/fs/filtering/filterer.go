package filtering

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Filterer defines the operations for narrowing a collection of filesystem items.
type Filterer interface {
	// Filter returns a new slice containing only the items that satisfy the filter criteria.
	Filter(items []shared.Item) ([]shared.Item, error)
}
