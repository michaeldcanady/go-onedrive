package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/core/environment"
	"github.com/michaeldcanady/go-onedrive/internal/core/shared"
)

// FileService is a persistent implementation of the profile.Service using JSON files.
type FileService struct {
	mu      sync.RWMutex
	baseDir string
}

// NewFileService initializes a new instance of the FileService.
func NewFileService(env environment.Service, baseDir string) (*FileService, error) {
	if baseDir == "" {
		configDir, err := env.ConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get config directory: %w", err)
		}
		baseDir = filepath.Join(configDir, "profiles")
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profile base directory: %w", err)
	}

	s := &FileService{
		baseDir: baseDir,
	}

	// Ensure default profile exists
	_, _ = s.Create(context.Background(), shared.DefaultProfileName)

	return s, nil
}

// Get returns the profile with the specified name if it exists on disk.
func (s *FileService) Get(ctx context.Context, name string) (shared.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.baseDir, name+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return shared.Profile{}, fmt.Errorf("profile %s not found: %w", name, err)
	}

	var p shared.Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return shared.Profile{}, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return p, nil
}

// List returns a list of all profiles available on disk.
func (s *FileService) List(ctx context.Context) ([]shared.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile directory: %w", err)
	}

	var profiles []shared.Profile
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".json")
		p, err := s.Get(ctx, name)
		if err == nil {
			profiles = append(profiles, p)
		}
	}

	return profiles, nil
}

// Create generates a new profile on disk with the specified name.
func (s *FileService) Create(ctx context.Context, name string) (shared.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.baseDir, name+".json")
	if _, err := os.Stat(filePath); err == nil {
		return shared.Profile{}, fmt.Errorf("profile %s already exists", name)
	}

	p := shared.Profile{
		Name: name,
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return shared.Profile{}, fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return shared.Profile{}, fmt.Errorf("failed to write profile file: %w", err)
	}

	return p, nil
}

// Update saves the specified profile to its corresponding JSON file on disk.
func (s *FileService) Update(ctx context.Context, p shared.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.baseDir, p.Name+".json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write profile file: %w", err)
	}

	return nil
}

// Delete removes the JSON file associated with the specified profile name from disk.
func (s *FileService) Delete(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if name == shared.DefaultProfileName {
		return fmt.Errorf("cannot delete the default profile")
	}

	filePath := filepath.Join(s.baseDir, name+".json")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete profile file: %w", err)
	}

	return nil
}

// Exists checks if a profile with the specified name exists on disk.
func (s *FileService) Exists(ctx context.Context, name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.baseDir, name+".json")
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
