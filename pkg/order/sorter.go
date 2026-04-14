package order

import (
	"slices"
	"sort"

	"github.com/michaeldcanady/go-onedrive/pkg/list"
)

// Sorter provides a multi-criteria sorting mechanism for a slice of items.
type Sorter[T any] struct {
	items       []T
	comparators list.LinkedList[Comparator[T]]
}

// NewSorter initializes a new Sorter with a clone of the given items.
func NewSorter[T any](items []T) *Sorter[T] {
	return &Sorter[T]{
		items: slices.Clone(items),
	}
}

// AddComparator appends a comparator to the end of the sorting criteria list.
func (s *Sorter[T]) AddComparator(cmp Comparator[T]) *Sorter[T] {
	s.comparators.PushBack(cmp)
	return s
}

// Sort orders the items according to the added comparators and returns the resulting slice.
func (s *Sorter[T]) Sort() []T {
	if s.comparators.Len() == 0 {
		return s.items
	}

	sort.SliceStable(s.items, func(i, j int) bool {
		for e := s.comparators.Front(); e != nil; e = e.Next() {
			cmp := e.Value()
			if cmp(s.items[i], s.items[j]) {
				return true
			}
			if cmp(s.items[j], s.items[i]) {
				return false
			}
		}
		return false
	})

	return s.items
}
