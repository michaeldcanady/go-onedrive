package filtering

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

var _ spec.Specification[shared.Item] = (*SizeFilter)(nil)

// SizeFilter is a specification that filters items based on their size.
type SizeFilter struct {
	minSize *int64
	maxSize *int64
}

// NewSizeFilter initializes a new SizeFilter with the provided min and max bounds.
func NewSizeFilter(minSize, maxSize *int64) *SizeFilter {
	return &SizeFilter{minSize: minSize, maxSize: maxSize}
}

// IsSatisfiedBy implements [spec.Specification].
func (s *SizeFilter) IsSatisfiedBy(candidate shared.Item) bool {
	if s.minSize != nil && candidate.Size < *s.minSize {
		return false
	}

	if s.maxSize != nil && candidate.Size > *s.maxSize {
		return false
	}

	return true
}
