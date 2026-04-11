package fs

import "time"

// Item provides metadata and identifying information for a filesystem object.
type Item struct {
	// ID is the unique identifier provided by the storage provider.
	ID string
	// Name is the display name of the item.
	Name string
	// Path is the absolute location of the item within the filesystem.
	Path string
	// Type indicates whether the item is a file or a folder.
	Type ItemType
	// Size is the length of the file's content in bytes.
	Size int64
	// ModifiedAt is the timestamp of the last modification.
	ModifiedAt time.Time
	// ETag is a unique version identifier for the item.
	ETag string
	// ProviderSpecific contains additional metadata unique to the storage provider.
	ProviderSpecific any
}

// FileTransferProgressEvent is emitted during file upload or download.
type FileTransferProgressEvent struct {
	Path        string
	Transferred int64
	Total       int64
}

func (e FileTransferProgressEvent) Name() string {
	return "fs.file_transfer_progress"
}
