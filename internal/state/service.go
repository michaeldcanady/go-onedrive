package state

// Service provides methods to retrieve, update, and manage the application state.
type Service interface {
	// Get retrieves a state value by its key.
	Get(key Key) (string, error)
	// Set assigns a value to a key within the specified scope.
	Set(key Key, value string, scope Scope) error
	// Clear removes a state value for the given key from all scopes.
	Clear(key Key) error

	// GetDriveAlias retrieves the drive ID associated with a user-defined alias.
	GetDriveAlias(alias string) (string, error)
	// SetDriveAlias assigns a drive ID to a specified alias.
	SetDriveAlias(alias, driveID string) error
	// RemoveDriveAlias deletes a drive alias.
	RemoveDriveAlias(alias string) error
	// ListDriveAliases returns all registered drive aliases.
	ListDriveAliases() (map[string]string, error)

	// GetDriveAliasByDriveID finds an alias associated with a given drive ID.
	GetDriveAliasByDriveID(driveID string) (string, error)
}
