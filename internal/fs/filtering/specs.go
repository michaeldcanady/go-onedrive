package filtering

import (
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

// HiddenSpec matches items that are considered hidden (start with a dot).
type HiddenSpec struct{}

func (s HiddenSpec) IsSatisfiedBy(item shared.Item) bool {
	return strings.HasPrefix(item.Name, ".")
}

// TypeSpec matches items of a specific filesystem type.
type TypeSpec struct {
	Type shared.ItemType
}

func (s TypeSpec) IsSatisfiedBy(item shared.Item) bool {
	if s.Type == shared.TypeUnknown {
		return true
	}
	return item.Type == s.Type
}

// NewFilterSpec creates a composite specification based on the provided filter options.
func NewFilterSpec(opts FilterOptions) spec.Specification[shared.Item] {
	var s spec.Specification[shared.Item] = spec.All[shared.Item]()

	// If not including all, exclude hidden files
	if !opts.IncludeAll {
		s = spec.And(s, spec.Not(HiddenSpec{}))
	}

	// Filter by item type if specified
	if opts.ItemType != shared.TypeUnknown {
		s = spec.And(s, TypeSpec{Type: opts.ItemType})
	}

	return s
}
