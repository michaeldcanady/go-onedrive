package shared

const (
	// DefaultProfileName is the name of the fallback profile.
	DefaultProfileName = "default"
)

// Scope determines the persistence level of state data.
type Scope int

const (
	// ScopeGlobal state persists across sessions.
	ScopeGlobal Scope = iota
	// ScopeSession state is transient and exists only for the duration of the current execution.
	ScopeSession
)
