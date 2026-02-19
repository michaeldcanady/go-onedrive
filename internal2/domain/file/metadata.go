package file

import "time"

type Metadata struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	FullPath   string     `json:"full_path"`
	Size       int64      `json:"size"`
	MimeType   string     `json:"mime_type"`
	ETag       string     `json:"etag"`
	CTag       string     `json:"ctag"`
	ParentID   string     `json:"parent_id"`
	CreatedAt  *time.Time `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at"`
	Type       ItemType   `json:"type"`
}
