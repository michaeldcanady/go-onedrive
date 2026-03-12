package shared

// CopyOptions defines the configuration for a copy operation.
type CopyOptions struct {
	// Recursive determines whether to copy directories and their contents.
	Recursive bool
	// Overwrite determines whether to replace an existing item at the destination.
	Overwrite bool
}
