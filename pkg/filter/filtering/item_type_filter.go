package filtering

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

var _ spec.Specification[shared.Item] = (*ItemTypeFilter)(nil)

// ItemTypeFilter is a specification that filters items based on their type.
type ItemTypeFilter struct {
	itemType shared.ItemType
}

// IsSatisfiedBy implements [spec.Specification].
func (i *ItemTypeFilter) IsSatisfiedBy(candidate shared.Item) bool {
	return candidate.Type == i.itemType
}
