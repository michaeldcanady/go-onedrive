package state

// Key identifies a piece of application state (e.g., active profile or drive).
type Key int

const (
	// KeyProfile represents the currently active profile.
	KeyProfile Key = iota
	// KeyDrive represents the currently active drive.
	KeyDrive
	// KeyAccessToken represents the cached authentication token.
	KeyAccessToken
)

// String returns the string representation of the Key.
func (k Key) String() string {
	switch k {
	case KeyProfile:
		return "profile"
	case KeyDrive:
		return "drive"
	case KeyAccessToken:
		return "access_token"
	default:
		return "unknown"
	}
}
