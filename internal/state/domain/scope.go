package domain

// Scope defines the persistence scope of a state change.
type Scope int

const (
	// ScopeGlobal indicates the change should be persisted to disk.
	ScopeGlobal Scope = iota
	// ScopeSession indicates the change should only last for the duration of the current process.
	ScopeSession
)
