package fs

import (
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/pkg/ignore"
)

// IgnoreMatcher defines the interface for checking if a path should be ignored.
type IgnoreMatcher interface {
	// ShouldIgnore returns true if the given relative path matches any ignore patterns.
	ShouldIgnore(path string, isDir bool) bool
}

// NoOpIgnoreMatcher is an implementation that ignores nothing.
type NoOpIgnoreMatcher struct{}

func (m NoOpIgnoreMatcher) ShouldIgnore(path string, isDir bool) bool {
	return false
}

// IgnoreMatcherFactory defines the interface for creating an IgnoreMatcher.
type IgnoreMatcherFactory interface {
	CreateMatcher(ctx context.Context, r io.Reader) (IgnoreMatcher, error)
}

// Ensure pkg/ignore implements the interface
var _ IgnoreMatcher = (*ignore.Matcher)(nil)
