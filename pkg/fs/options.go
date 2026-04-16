package fs

// CopyOptions defines the configuration for a copy operation.
type CopyOptions struct {
	// Recursive determines whether to include nested items.
	Recursive bool
	// Overwrite determines whether to replace existing items at the destination.
	Overwrite bool
}

// ListOptions defines the configuration for an enumeration operation.
type ListOptions struct {
	// Recursive determines whether to traverse into subdirectories.
	Recursive bool
}

// ReadOptions defines the configuration for a file reading operation.
type ReadOptions struct{}

// WriteOptions defines the configuration for a write operation.
type WriteOptions struct {
	// Overwrite determines whether to replace an existing item.
	Overwrite bool
	// IfMatch is the ETag of the item that should be overwritten.
	IfMatch string
	// Size is the total number of bytes to be written (required for resumable uploads).
	Size int64
}
