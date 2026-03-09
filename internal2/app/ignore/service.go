package ignore

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	domain "github.com/michaeldcanady/go-onedrive/internal2/domain/ignore"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/pkg/ignore"
)

var _ domain.Service = (*service)(nil)

type service struct {
	mu            sync.RWMutex
	globalMatcher *ignore.Matcher
	registry      map[string]*ignore.Matcher
	logger        logging.Logger
}

// NewService creates a new implementation of the domain.Service.
func NewService(logger logging.Logger) domain.Service {
	return &service{
		logger:        logger,
		globalMatcher: ignore.NewMatcher(nil),
		registry:      make(map[string]*ignore.Matcher),
	}
}

func (s *service) ShouldIgnore(ctx context.Context, path string, isDir bool) (bool, *ignore.Rule) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 1. Check global rules (CLI flags, etc.)
	if ignored, rule := s.globalMatcher.MatchWithRule(path, isDir); ignored {
		return true, rule
	}

	// Clean path for registry lookups
	path = filepath.Clean(path)
	segments := strings.Split(path, string(os.PathSeparator))

	// 2. Check hierarchical rules
	// We check matchers from the root down to the file's directory.
	// Note: rules in subdirectories can negate rules from parent directories.
	
	currentPath := "."
	var lastIgnored bool
	var lastRule *ignore.Rule

	// Check root directory matcher
	if m, ok := s.registry[currentPath]; ok {
		if ignored, rule := m.MatchWithRule(path, isDir); ignored {
			lastIgnored = true
			lastRule = rule
		} else if rule != nil { // Matched but negated
			lastIgnored = false
			lastRule = rule
		}
	}

	// Check subdirectory matchers
	accPath := ""
	for i := 0; i < len(segments)-1; i++ {
		if accPath == "" {
			accPath = segments[i]
		} else {
			accPath = filepath.Join(accPath, segments[i])
		}

		if m, ok := s.registry[accPath]; ok {
			// Calculate path relative to this ignore file's directory
			relPath, _ := filepath.Rel(accPath, path)
			if ignored, rule := m.MatchWithRule(relPath, isDir); ignored {
				lastIgnored = true
				lastRule = rule
			} else if rule != nil {
				lastIgnored = false
				lastRule = rule
			}
		}
	}

	return lastIgnored, lastRule
}

func (s *service) LoadGlobalPatterns(ctx context.Context, patterns []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var newRules []ignore.Rule
	for _, p := range patterns {
		l := ignore.NewLexer(p)
		parser := ignore.NewParser(l)
		doc := parser.Parse()
		newRules = append(newRules, doc.Rules...)
	}

	s.globalMatcher = ignore.NewMatcher(append(s.globalMatcher.Rules(), newRules...))
	return nil
}

func (s *service) LoadIgnoreFile(ctx context.Context, path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Debug("ignore file not found, skipping", logging.String("path", path))
			return nil
		}
		return err
	}
	defer f.Close()

	matcher, err := ignore.CompileReader(f)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	s.registry[dir] = matcher
	
	s.logger.Debug("loaded ignore file", 
		logging.String("path", path), 
		logging.String("scope", dir),
		logging.Int("rules", len(matcher.Rules())))
		
	return nil
}
