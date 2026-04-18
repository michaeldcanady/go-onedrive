package filtering

import (
	"time"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/spec"
)

var _ spec.Specification[shared.Item] = (*DateFilter)(nil)

// DateFilter is a specification that filters items based on their modification date.
type DateFilter struct {
	before *time.Time
	after  *time.Time
}

// NewDateFilter initializes a new DateFilter with the provided before and after timestamps.
func NewDateFilter(before, after *time.Time) *DateFilter {
	return &DateFilter{before: before, after: after}
}

// IsSatisfiedBy implements [spec.Specification].
func (d *DateFilter) IsSatisfiedBy(candidate shared.Item) bool {
	if d.before != nil && !candidate.ModifiedAt.Before(*d.before) {
		return false
	}

	if d.after != nil && !candidate.ModifiedAt.After(*d.after) {
		return false
	}

	return true
}
