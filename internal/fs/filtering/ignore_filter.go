package filtering

import (
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/ignore"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

var _ spec.Specification[shared.Item] = (*IgnoreFilter)(nil)

// IgnoreFilter is a specification that filters items based on ignore patterns.
type IgnoreFilter struct {
	matcher *ignore.Matcher
}

// NewIgnoreFilter initializes a new IgnoreFilter with the provided matcher.
func NewIgnoreFilter(matcher *ignore.Matcher) *IgnoreFilter {
	return &IgnoreFilter{matcher: matcher}
}

// IsSatisfiedBy implements [spec.Specification].
func (i *IgnoreFilter) IsSatisfiedBy(candidate shared.Item) bool {
	if i.matcher == nil {
		return true
	}

	isDir := candidate.Type == shared.TypeFolder
	// ShouldIgnore returns true if the item matches an ignore pattern.
	// We want to return true (satisfied) if the item should NOT be ignored.
	return !i.matcher.ShouldIgnore(candidate.Path, isDir)
}
