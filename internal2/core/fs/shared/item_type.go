package shared

// ItemType represents the classification of a filesystem entry.
type ItemType int

const (
	// TypeFile identifies the item as a regular file.
	TypeFile ItemType = iota
	// TypeFolder identifies the item as a directory or folder.
	TypeFolder
)
