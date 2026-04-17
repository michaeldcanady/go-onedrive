package state

// Service provides methods to retrieve, update, and manage the application state.
type Service interface {
	// Get retrieves a state value by its key.
	Get(key Key) (string, error)
	// Set assigns a value to a key within the specified scope.
	Set(key Key, value string, scope Scope) error
	// Clear removes a state value for the given key from all scopes.
	Clear(key Key) error

	// GetScoped retrieves a value from a named sub-bucket (e.g., "tokens").
	GetScoped(bucket, key string) (string, error)
	// SetScoped assigns a value to a key within a named sub-bucket.
	SetScoped(bucket, key, value string, scope Scope) error
	// ClearScoped removes a value from a named sub-bucket.
	ClearScoped(bucket, key string) error
	// ListScoped returns all keys within a named sub-bucket across all scopes.
	ListScoped(bucket string) ([]string, error)
}
