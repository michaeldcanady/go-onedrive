package state

// Key identifies a piece of application state.
// State keys are used for transient, session-scoped data.
type Key int

const (
	// KeyProfile represents the currently active profile (transient override).
	KeyProfile Key = iota
	// KeyConfigOverride represents a transient configuration path override.
	KeyConfigOverride
)

// String returns the string representation of the Key.
func (k Key) String() string {
	switch k {
	case KeyProfile:
		return "profile"
	case KeyConfigOverride:
		return "config_override"
	default:
		return "unknown"
	}
}
