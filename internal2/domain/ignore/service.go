package ignore

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/pkg/ignore"
)

// Service defines the domain-level contract for the ignore system.
// It orchestrates the loading, merging, and evaluation of ignore patterns
// using the underlying pkg/ignore engine.
type Service interface {
	// ShouldIgnore evaluates if a specific path matches any active ignore rules.
	// It checks rules in hierarchical order (e.g., global -> root -> subdirectories).
	ShouldIgnore(ctx context.Context, path string, isDir bool) (bool, *ignore.Rule)

	// LoadGlobalPatterns adds patterns that apply to all paths (e.g., from CLI flags).
	LoadGlobalPatterns(ctx context.Context, patterns []string) error

	// LoadIgnoreFile reads patterns from a file and associates them with that file's directory.
	// Rules in this file will only apply to the directory and its descendants.
	LoadIgnoreFile(ctx context.Context, path string) error
}
