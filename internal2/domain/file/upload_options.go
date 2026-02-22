package file

type UploadOptions struct {
	// If true, ignore ETag mismatch errors and force overwrite
	Force bool

	// If true, do not update the content cache after upload
	NoStore bool
}
