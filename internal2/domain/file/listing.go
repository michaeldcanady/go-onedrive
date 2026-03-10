package file

// Listing represents a cached collection of drive item IDs at a specific path.
type Listing struct {
	// ETag is the entity tag of the parent folder when the listing was fetched.
	ETag string
	// ChildIDs is a slice of IDs for the items contained within the folder.
	ChildIDs []string
}
