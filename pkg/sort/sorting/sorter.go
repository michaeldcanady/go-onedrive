package sorting

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
)

// Sorter defines the operations for ordering a collection of filesystem items.
type Sorter interface {
	// Sort returns a new slice with items ordered according to the sorter's internal logic.
	Sort(items []shared.Item) ([]shared.Item, error)
}
