package shared

// WriteOptions defines the configuration for a write operation.
type WriteOptions struct {
	// Overwrite determines whether to replace an existing item.
	Overwrite bool
	// IfMatch is the ETag of the item that should be overwritten.
	IfMatch string
}
