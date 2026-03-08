package fs

// IgnoreMatcher defines the interface for determining if a path should be ignored.
type IgnoreMatcher interface {
	// ShouldIgnore returns true if the given relative path matches any ignore patterns.
	ShouldIgnore(path string, isDir bool) bool
}

// NoOpIgnoreMatcher is an implementation that ignores nothing.
type NoOpIgnoreMatcher struct{}

func (m NoOpIgnoreMatcher) ShouldIgnore(path string, isDir bool) bool {
	return false
}
