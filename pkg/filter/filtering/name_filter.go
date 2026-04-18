package filtering

import (
	"path/filepath"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

var _ spec.Specification[shared.Item] = (*NameFilter)(nil)

// NameFilter is a specification that filters items based on their name matching one or more glob patterns.
type NameFilter struct {
	patterns []string
}

// NewNameFilter initializes a new NameFilter with the provided glob patterns.
func NewNameFilter(patterns []string) *NameFilter {
	return &NameFilter{patterns: patterns}
}

// IsSatisfiedBy implements [spec.Specification].
func (n *NameFilter) IsSatisfiedBy(candidate shared.Item) bool {
	if len(n.patterns) == 0 {
		return true
	}

	for _, pattern := range n.patterns {
		matched, err := filepath.Match(pattern, candidate.Name)
		if err == nil && matched {
			return true
		}
	}

	return false
}
