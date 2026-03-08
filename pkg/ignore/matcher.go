package ignore

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"
)

// Matcher evaluates paths against patterns.
type Matcher struct {
	patterns []Pattern
}

// NewMatcher creates a Matcher from a slice of patterns.
func NewMatcher(patterns []Pattern) *Matcher {
	return &Matcher{patterns: patterns}
}

// ParseReader reads from an io.Reader and returns a Matcher.
func ParseReader(r io.Reader) (*Matcher, error) {
	var sb strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	l := NewLexer(sb.String())
	p := NewParser(l)
	return NewMatcher(p.Parse()), nil
}

// ShouldIgnore checks if a path matches the patterns.
func (m *Matcher) ShouldIgnore(path string, isDir bool) bool {
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimPrefix(path, "/")

	// Split path into segments to check each parent
	segments := strings.Split(path, string(charSlash))
	
	ignored := false
	currentPath := ""

	for i, segment := range segments {
		if currentPath == "" {
			currentPath = segment
		} else {
			currentPath += string(charSlash) + segment
		}

		isCurrentDir := isDir
		if i < len(segments)-1 {
			isCurrentDir = true
		}

		for _, p := range m.patterns {
			if p.IsDir && !isCurrentDir {
				continue
			}

			if matchPattern(p.Path, currentPath) {
				ignored = !p.IsNegate
			}
		}
	}

	return ignored
}

func matchPattern(pattern, path string) bool {
	// 1. Exact match
	if pattern == path {
		return true
	}

	// 2. Directory prefix match (e.g., "node_modules" matches "node_modules/foo")
	if strings.HasPrefix(path, pattern+string(charSlash)) {
		return true
	}

	// 3. Simple wildcard matching via filepath.Match
	match, _ := filepathMatch(pattern, path)
	return match
}

// filepathMatch is a placeholder or wrapper for actual matching logic.
func filepathMatch(pattern, path string) (bool, error) {
	parts := strings.Split(path, string(charSlash))
	for _, part := range parts {
		if ok, _ := filepath.Match(pattern, part); ok {
			return true, nil
		}
	}
	return false, nil
}
