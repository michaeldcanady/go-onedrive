package sorting

import (
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/order"
)

// ItemComparator defines a function for comparing two filesystem items.
type ItemComparator order.Comparator[shared.Item]

// Reverse returns an ItemComparator that reverses the order of the given comparator.
func Reverse(cmp ItemComparator) ItemComparator {
	return func(i, j shared.Item) bool {
		return cmp(j, i)
	}
}

// CompareByName compares two items by their name.
func CompareByName(i, j shared.Item) bool {
	return strings.ToLower(i.Name) < strings.ToLower(j.Name)
}

// CompareBySize compares two items by their size.
func CompareBySize(i, j shared.Item) bool {
	return i.Size < j.Size
}

// CompareByPath compares two items by their path.
func CompareByPath(i, j shared.Item) bool {
	return strings.ToLower(i.Path) < strings.ToLower(j.Path)
}

// CompareByType compares two items by their type.
func CompareByType(i, j shared.Item) bool {
	return i.Type < j.Type
}

// CompareByModifiedAt compares two items by their last modification time.
func CompareByModifiedAt(i, j shared.Item) bool {
	return i.ModifiedAt.Before(j.ModifiedAt)
}
