package shared

// ListOptions provides filtering and traversal settings for directory listing.
type ListOptions struct {
	// Recursive determines whether to include nested children of directories.
	Recursive bool
}
