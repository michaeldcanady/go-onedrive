package fs

type ListOptions struct {
	Recursive bool
	// SkipCache indicates to query API even if cache is current
	SkipCache bool
	// NoCache indicates not to cache the results
	NoCache bool
}
