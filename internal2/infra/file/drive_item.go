package file

import "time"

type DriveItem struct {
	DriveID          string
	ID               string
	Name             string
	Path             string
	PathWithoutDrive string
	IsFolder         bool
	Size             int64
	ETag             string
	MimeType         string
	Modified         time.Time
}
