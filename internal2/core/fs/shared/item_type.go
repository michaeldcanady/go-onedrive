package shared

// ItemType represents the classification of a filesystem entry.
type ItemType int

const (
	// TypeUnknown represents an unclassified filesystem item.
	TypeUnknown ItemType = iota
	// TypeFile identifies the item as a regular file.
	TypeFile
	// TypeFolder identifies the item as a directory or folder.
	TypeFolder
)
