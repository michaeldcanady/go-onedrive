package drive

// Drive represents the metadata for a OneDrive or local drive.
type Drive struct {
	// ID is the unique identifier for the drive.
	ID string
	// Name is the display name of the drive.
	Name string
	// Type specifies the category of the drive (personal, business, etc.).
	Type DriveType
	// Owner is the name or identifier of the drive's owner.
	Owner string
	// ReadOnly indicates if the drive is read-only for the current user.
	ReadOnly bool
}
