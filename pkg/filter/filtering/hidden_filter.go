package filtering

import (
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

var _ spec.Specification[shared.Item] = (*HiddenFilter)(nil)

// HiddenFilter is a specification that filters items based on their hidden status.
type HiddenFilter struct {
	hidden bool
}

// IsSatisfiedBy implements [spec.Specification].
func (h *HiddenFilter) IsSatisfiedBy(candidate shared.Item) bool {
	return strings.HasPrefix(candidate.Name, ".") == h.hidden
}
