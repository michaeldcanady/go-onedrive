package drive

// DriveType represents the category of a drive (e.g., personal, business, SharePoint).
type DriveType int

const (
	// DriveTypeUnknown represents an unknown drive type.
	DriveTypeUnknown DriveType = iota
	// DriveTypePersonal represents a personal OneDrive drive.
	DriveTypePersonal
	// DriveTypeBusiness represents a business OneDrive drive.
	DriveTypeBusiness
	// DriveTypeSharePoint represents a SharePoint document library drive.
	DriveTypeSharePoint
)

// String returns the string representation of the drive type.
func (dt DriveType) String() string {
	switch dt {
	case DriveTypePersonal:
		return "personal"
	case DriveTypeBusiness:
		return "business"
	case DriveTypeSharePoint:
		return "sharepoint"
	default:
		return "unknown"
	}
}

// NewDriveType converts a string to its corresponding DriveType.
func NewDriveType(s string) DriveType {
	switch s {
	case "personal":
		return DriveTypePersonal
	case "business":
		return DriveTypeBusiness
	case "sharepoint":
		return DriveTypeSharePoint
	default:
		return DriveTypeUnknown
	}
}
