package state

// Scope determines the persistence level of state data.
type Scope int

const (
	// ScopeGlobal state persists across sessions.
	ScopeGlobal Scope = iota
	// ScopeSession state is transient and exists only for the duration of the current execution.
	ScopeSession
)

// String returns a string representation of the scope.
func (s Scope) String() string {
	switch s {
	case ScopeGlobal:
		return "global"
	case ScopeSession:
		return "session"
	default:
		return "unknown"
	}
}
