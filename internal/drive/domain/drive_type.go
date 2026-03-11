package domain

// DriveType represents the type of drive (e.g., personal, business, SharePoint).
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

// NewDriveType creates a new DriveType from a string representation.
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
