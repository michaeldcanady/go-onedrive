package fs

// URI represents a structured filesystem location.
type URI struct {
	// Provider is the name of the filesystem provider (e.g., "local", "onedrive").
	Provider string
	// DriveID is the optional identifier for a specific drive within the provider.
	DriveID string
	// Path is the location of the item within the provider or drive.
	Path string
}

// String returns the canonical string representation of the URI.
func (u *URI) String() string {
	if u.DriveID != "" {
		return u.Provider + ":" + u.DriveID + ":" + u.Path
	}
	return u.Provider + ":" + u.Path
}
