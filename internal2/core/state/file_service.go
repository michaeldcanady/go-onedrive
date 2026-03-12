package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/core/shared"
)

// FileService is a persistent implementation of the state.Service using a JSON file.
type FileService struct {
	mu       sync.RWMutex
	filePath string
	global   map[string]string
	session  map[string]string
}

// NewFileService initializes a new instance of the FileService.
func NewFileService(filePath string) (*FileService, error) {
	if filePath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user config directory: %w", err)
		}
		filePath = filepath.Join(configDir, "odc", "state.json")
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	s := &FileService{
		filePath: filePath,
		global:   make(map[string]string),
		session:  make(map[string]string),
	}

	if err := s.load(); err != nil {
		// If the file doesn't exist, we just start with an empty state.
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Set default profile if not set
	if _, err := s.Get(KeyProfile); err != nil {
		_ = s.Set(KeyProfile, shared.DefaultProfileName, ScopeGlobal)
	}

	return s, nil
}

// Get retrieves a state value by its key, checking session scope first, then global.
func (s *FileService) Get(key Key) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keyStr := key.String()

	if val, ok := s.session[keyStr]; ok {
		return val, nil
	}

	if val, ok := s.global[keyStr]; ok {
		return val, nil
	}

	return "", fmt.Errorf("state key %s not found", keyStr)
}

// Set assigns a value to a key within the specified scope.
func (s *FileService) Set(key Key, value string, scope Scope) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyStr := key.String()

	switch scope {
	case ScopeSession:
		s.session[keyStr] = value
		return nil
	case ScopeGlobal:
		s.global[keyStr] = value
		return s.save()
	default:
		return fmt.Errorf("unsupported scope: %v", scope)
	}
}

// Clear removes a state value for the given key from all scopes.
func (s *FileService) Clear(key Key) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyStr := key.String()
	delete(s.session, keyStr)
	delete(s.global, keyStr)

	return s.save()
}

// GetDriveAlias retrieves the drive ID associated with a user-defined alias.
func (s *FileService) GetDriveAlias(alias string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.global["drive_alias:"+alias]; ok {
		return val, nil
	}
	return "", fmt.Errorf("drive alias %s not found", alias)
}

// SetDriveAlias assigns a drive ID to a specified alias.
func (s *FileService) SetDriveAlias(alias, driveID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.global["drive_alias:"+alias] = driveID
	return s.save()
}

// RemoveDriveAlias deletes a drive alias.
func (s *FileService) RemoveDriveAlias(alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.global, "drive_alias:"+alias)
	return s.save()
}

// ListDriveAliases returns all registered drive aliases.
func (s *FileService) ListDriveAliases() (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	aliases := make(map[string]string)
	prefix := "drive_alias:"
	for k, v := range s.global {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			aliases[k[len(prefix):]] = v
		}
	}
	return aliases, nil
}

func (s *FileService) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &s.global); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return nil
}

func (s *FileService) save() error {
	data, err := json.MarshalIndent(s.global, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}
