package config

import (
	"embed"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"gopkg.in/yaml.v3"
)

//go:embed defaults.yaml
var defaultFS embed.FS

type configService struct {
	repo     Repository
	logger   logger.Service
	defaults map[string]any
}

// NewConfigService returns a new [Service] initialized with the provided repository.
// It pre-populates default settings from an embedded defaults.yaml file.
func NewConfigService(repo Repository, l logger.Service) Service {
	s := &configService{
		repo:   repo,
		logger: l,
	}

	defaults := make(map[string]any)
	data, err := defaultFS.ReadFile("defaults.yaml")
	if err != nil {
		l.Error("failed to read embedded defaults.yaml", "error", err)
	} else if err := yaml.Unmarshal(data, &defaults); err != nil {
		l.Error("failed to unmarshal embedded defaults.yaml", "error", err)
	}

	s.defaults = defaults
	return s
}

func (s *configService) Get(key string) (any, error) {
	config, err := s.repo.Load()
	if err != nil {
		return nil, err
	}

	val, _, err := s.traverse(config, key, false)
	if err != nil || val == nil {
		// Fallback to defaults
		val, _, err = s.traverse(s.defaults, key, false)
		if err != nil || val == nil {
			return nil, fmt.Errorf("key not found: %s", key)
		}
	}

	return val, nil
}

func (s *configService) Set(key string, value string) error {
	config, err := s.repo.Load()
	if err != nil {
		return err
	}

	_, _, err = s.traverse(config, key, true)
	if err != nil {
		return err
	}

	// Traversal only gets us to the map containing the key, we still need to set it
	// Actually, let's refine traverse to return the leaf map and the last part of the key
	m, lastPart, err := s.traverse(config, key, true)
	if err != nil {
		return err
	}
	m[lastPart] = value

	return s.repo.Save(config)
}

func (s *configService) All() (map[string]any, error) {
	return s.repo.Load()
}

func (s *configService) traverse(config map[string]any, key string, createMissing bool) (map[string]any, string, error) {
	parts := strings.Split(key, ".")
	current := config

	for i, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]any)
		if !ok {
			if createMissing {
				next = make(map[string]any)
				current[part] = next
			} else {
				return nil, "", fmt.Errorf("key path not found: %s", strings.Join(parts[:i+1], "."))
			}
		}
		current = next
	}

	return current, parts[len(parts)-1], nil
}

func (s *configService) getValue(config map[string]any, key string) any {
	m, lastPart, err := s.traverse(config, key, false)
	if err != nil {
		return nil
	}
	return m[lastPart]
}

func (s *configService) setValue(config map[string]any, key string, value any) {
	m, lastPart, err := s.traverse(config, key, true)
	if err != nil {
		return
	}
	m[lastPart] = value
}
