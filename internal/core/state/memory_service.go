package state

import "fmt"

// MemoryService is an in-memory implementation of the Service interface.
// State stored here is transient and does not persist after the process exits.
type MemoryService struct {
	// global stores state values that represent cross-session data.
	global map[Key]string
	// session stores state values that are transient for the current execution.
	session map[Key]string
	// aliases stores user-defined names for specific drive IDs.
	aliases map[string]string
}

// NewMemoryService initializes a new instance of MemoryService.
func NewMemoryService() *MemoryService {
	return &MemoryService{
		global:  make(map[Key]string),
		session: make(map[Key]string),
		aliases: make(map[string]string),
	}
}

// Get retrieves a value for the specified key, prioritizing session scope.
func (s *MemoryService) Get(key Key) (string, error) {
	if val, ok := s.session[key]; ok {
		return val, nil
	}
	if val, ok := s.global[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("state key %s not found", key)
}

// Set assigns a value to a key within the specified scope.
func (s *MemoryService) Set(key Key, value string, scope Scope) error {
	if scope == ScopeSession {
		s.session[key] = value
	} else {
		s.global[key] = value
	}
	return nil
}

// Clear removes a key and its value from both session and global scopes.
func (s *MemoryService) Clear(key Key) error {
	delete(s.session, key)
	delete(s.global, key)
	return nil
}

// GetDriveAlias retrieves the drive ID mapped to the given alias.
func (s *MemoryService) GetDriveAlias(alias string) (string, error) {
	if val, ok := s.aliases[alias]; ok {
		return val, nil
	}
	return "", fmt.Errorf("alias %s not found", alias)
}

// SetDriveAlias assigns a drive ID to the specified alias name.
func (s *MemoryService) SetDriveAlias(alias, driveID string) error {
	s.aliases[alias] = driveID
	return nil
}

// RemoveDriveAlias deletes the specified alias mapping.
func (s *MemoryService) RemoveDriveAlias(alias string) error {
	delete(s.aliases, alias)
	return nil
}

// ListDriveAliases returns a map of all registered drive aliases.
func (s *MemoryService) ListDriveAliases() (map[string]string, error) {
	return s.aliases, nil
}
