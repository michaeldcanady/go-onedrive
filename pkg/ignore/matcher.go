package ignore

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"
)

// Matcher evaluates paths against rules.
type Matcher struct {
	rules []Rule
}

// NewMatcher creates a new Matcher.
func NewMatcher(rules []Rule) *Matcher {
	return &Matcher{rules: rules}
}

// Compile parses a raw ignore string and returns a Matcher.
func Compile(input string) *Matcher {
	l := NewLexer(input)
	p := NewParser(l)
	doc := p.Parse()
	return NewMatcher(doc.Rules)
}

// CompileReader reads from an io.Reader and returns a Matcher.
func CompileReader(r io.Reader) (*Matcher, error) {
	var sb strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteByte(charNewline)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return Compile(sb.String()), nil
}

// Match checks if a path is ignored.
func (m *Matcher) Match(path string, isDir bool) bool {
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimPrefix(path, string(charSlash))
	
	// Split path into segments
	segments := strings.Split(path, string(charSlash))
	ignored := false
	
	currentPath := ""
	for i, segment := range segments {
		if currentPath == "" {
			currentPath = segment
		} else {
			currentPath += string(charSlash) + segment
		}
		
		isSegmentDir := true
		if i == len(segments)-1 {
			isSegmentDir = isDir
		} else {
			isSegmentDir = true
		}
		
		for _, rule := range m.rules {
			// Optimization: if rule is directory-only but this segment is not a dir, skip
			if rule.DirOnly && !isSegmentDir {
				continue
			}
			
			matched := matchPattern(rule.Pattern, currentPath, segment)
			if matched {
				if rule.Negate {
					ignored = false
				} else {
					ignored = true
				}
			}
		}
	}
	
	return ignored
}

func matchPattern(pattern, fullPath, segment string) bool {
	// 1. If pattern has no slash (and not anchored), match against segment (basename).
	// Exception: `**/foo` -> treated as "match foo anywhere" effectively.
	
	// If pattern starts with `/`, it is anchored.
	isAnchored := strings.HasPrefix(pattern, string(charSlash))
	pattern = strings.TrimPrefix(pattern, string(charSlash))
	
	hasSlash := strings.Contains(pattern, string(charSlash))
	
	if !hasSlash && !isAnchored {
		// Match against basename
		ok, _ := filepath.Match(pattern, segment)
		return ok
	}
	
	// Anchored or has slash: match against full path
	ok, _ := filepath.Match(pattern, fullPath)
	return ok
}
