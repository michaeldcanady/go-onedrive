package file

type UploadOptions struct {
	// If true, ignore ETag mismatch errors and force overwrite
	Force bool

	// If set, the upload will only proceed if the current ETag matches
	IfMatch string

	// If true, do not update the content cache after upload
	NoStore bool
}
