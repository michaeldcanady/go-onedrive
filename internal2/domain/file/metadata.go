package file

import "time"

type Metadata struct {
	ID         string
	Name       string
	Path       string
	Size       int64
	MimeType   string
	ETag       string
	CTag       string
	ParentID   string
	CreatedAt  *time.Time
	ModifiedAt *time.Time
	Type       ItemType
}
