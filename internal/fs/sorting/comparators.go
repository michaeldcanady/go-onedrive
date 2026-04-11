package sorting

import (
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/sortutil"
)

// NameComparator sorts items alphabetically by their name.
type NameComparator struct{}

func (c NameComparator) Less(i, j shared.Item) bool {
	return strings.ToLower(i.Name) < strings.ToLower(j.Name)
}

// SizeComparator sorts items by their byte size.
type SizeComparator struct{}

func (c SizeComparator) Less(i, j shared.Item) bool {
	return i.Size < j.Size
}

// DateComparator sorts items by their last modified timestamp.
type DateComparator struct{}

func (c DateComparator) Less(i, j shared.Item) bool {
	return i.ModifiedAt.Before(j.ModifiedAt)
}

// FolderFirstComparator ensures folders always appear before files.
type FolderFirstComparator struct{}

func (c FolderFirstComparator) Less(i, j shared.Item) bool {
	if i.Type == shared.TypeFolder && j.Type != shared.TypeFolder {
		return true
	}
	return false
}

// NewItemComparator builds a composite comparator based on sorting options.
func NewItemComparator(opts SortingOptions) sortutil.Comparator[shared.Item] {
	var primary sortutil.Comparator[shared.Item]

	// Default to folders first if not explicitly disabled (common coreutils behavior)
	primary = FolderFirstComparator{}

	var secondary sortutil.Comparator[shared.Item]
	switch strings.ToLower(opts.Field) {
	case "size":
		secondary = SizeComparator{}
	case "modifiedat", "date", "time":
		secondary = DateComparator{}
	default: // Default to name
		secondary = NameComparator{}
	}

	// Apply direction to the field-specific comparator
	if opts.Direction == DirectionDescending {
		secondary = sortutil.Reverse(secondary)
	}

	return sortutil.Then(primary, secondary)
}
