package fs

// ItemType distinguishes between different categories of filesystem objects.
type ItemType int

const (
	// TypeUnknown identifies an item with an unresolved or invalid category.
	TypeUnknown ItemType = iota
	// TypeFile identifies an object containing data (e.g., a document or image).
	TypeFile
	// TypeFolder identifies an object that can contain other items.
	TypeFolder
)
