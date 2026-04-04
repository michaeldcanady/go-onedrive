package state

// Service provides methods to retrieve, update, and manage the application state.
type Service interface {
	// Get retrieves a state value by its key.
	Get(key Key) (string, error)
	// Set assigns a value to a key within the specified scope.
	Set(key Key, value string, scope Scope) error
	// Clear removes a state value for the given key from all scopes.
	Clear(key Key) error
}
