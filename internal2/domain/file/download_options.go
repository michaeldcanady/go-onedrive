package file

type DownloadOptions struct {
	// If true, bypass cache entirely and always fetch from OneDrive
	NoCache bool

	// If true, do NOT write the downloaded content into the cache
	NoStore bool
}
