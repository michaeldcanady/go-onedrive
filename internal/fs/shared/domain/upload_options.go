package domain

type UploadOptions struct {
	Overwrite bool
	IfMatch   string
	NoStore   bool
}
