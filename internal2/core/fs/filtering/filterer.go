package filtering

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
)

// Filterer defines the operations for narrowing a collection of filesystem items.
type Filterer interface {
	// Filter returns a new slice containing only the items that satisfy the filter criteria.
	Filter(items []shared.Item) ([]shared.Item, error)
}
