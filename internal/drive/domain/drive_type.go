package domain

// TODO: Migrate to int-based enum
type DriveType string

const (
	DriveTypePersonal   DriveType = "personal"
	DriveTypeBusiness   DriveType = "business"
	DriveTypeSharePoint DriveType = "sharepoint"
)
