package profile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/core/shared"
)

// FileService is a filesystem-based implementation of the Profile Service.
type FileService struct {
	mu      sync.RWMutex
	baseDir string
}

// NewFileService initializes a new instance of FileService.
func NewFileService(baseDir string) (*FileService, error) {
	if baseDir == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user config directory: %w", err)
		}
		baseDir = filepath.Join(configDir, "odc", "profiles")
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profile base directory: %w", err)
	}

	s := &FileService{
		baseDir: baseDir,
	}

	// Ensure default profile exists
	if _, err := s.Create(context.Background(), shared.DefaultProfileName); err != nil {
		// If it already exists, that's fine. We don't need to return an error here.
		// A real implementation might want to log this or check if it's a different error.
	}

	return s, nil
}

// Get returns the profile with the specified name if it exists on disk.
func (s *FileService) Get(ctx context.Context, name string) (shared.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.baseDir, name)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return shared.Profile{}, fmt.Errorf("profile %s not found", name)
		}
		return shared.Profile{}, fmt.Errorf("failed to access profile: %w", err)
	}

	if !info.IsDir() {
		return shared.Profile{}, fmt.Errorf("profile path %s is not a directory", path)
	}

	return shared.Profile{
		Name:       name,
		Path:       path,
		ConfigPath: filepath.Join(path, "config.yaml"),
	}, nil
}

// List returns all profiles found in the base directory.
func (s *FileService) List(ctx context.Context) ([]shared.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile directory: %w", err)
	}

	var list []shared.Profile
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(s.baseDir, entry.Name())
			list = append(list, shared.Profile{
				Name:       entry.Name(),
				Path:       path,
				ConfigPath: filepath.Join(path, "config.yaml"),
			})
		}
	}

	return list, nil
}

// Create creates a new profile directory on disk.
func (s *FileService) Create(ctx context.Context, name string) (shared.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.baseDir, name)
	if _, err := os.Stat(path); err == nil {
		// Profile already exists
		return shared.Profile{
			Name:       name,
			Path:       path,
			ConfigPath: filepath.Join(path, "config.yaml"),
		}, nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return shared.Profile{}, fmt.Errorf("failed to create profile directory: %w", err)
	}

	return shared.Profile{
		Name:       name,
		Path:       path,
		ConfigPath: filepath.Join(path, "config.yaml"),
	}, nil
}

// Delete removes a profile directory and all its contents.
func (s *FileService) Delete(ctx context.Context, name string) error {
	if name == shared.DefaultProfileName {
		return fmt.Errorf("cannot delete the default profile")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.baseDir, name)
	return os.RemoveAll(path)
}

// Exists checks if a profile directory exists on disk.
func (s *FileService) Exists(ctx context.Context, name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.baseDir, name)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check profile existence: %w", err)
	}

	return info.IsDir(), nil
}
