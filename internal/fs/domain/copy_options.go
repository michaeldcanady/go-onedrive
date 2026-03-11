package domain

// CopyOptions defines options for the Copy operation.
type CopyOptions struct {
	// Overwrite specifies whether to overwrite the destination if it exists.
	Overwrite bool
	// Recursive specifies whether to copy folders recursively.
	Recursive bool
	// Matcher specifies a matcher to use for ignoring items during a recursive copy.
	Matcher IgnoreMatcher
}
