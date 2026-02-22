package file

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// mapItemToMetadata converts a Microsoft Graph DriveItem into the internal
// file.Metadata domain model.
//
// The function extracts:
//
//   - Parent ID and normalized parent path.
//   - File vs folder type (via presence of File facet).
//   - MIME type (files only).
//   - Size, ETag, CTag.
//   - Created/modified timestamps.
//
// Path normalization:
//
//   - Graph parent paths are of the form: "driveID:/path/to/folder".
//   - The function strips the "{driveID}:" prefix.
//   - Trailing slashes are removed.
//
// The returned Metadata is always non‑nil and contains zero values for any
// missing Graph fields.
func mapItemToMetadata(it models.DriveItemable) *file.Metadata {
	var (
		parentID string
		mimeType string
		fullPath string
		path     string
		ype      file.ItemType = file.ItemTypeFolder
	)
	if parent := it.GetParentReference(); parent != nil {
		parentID = deref(parent.GetId())
		fullPath = deref(parent.GetPath())
		if fullPath != "" {
			path = strings.Split(fullPath, ":")[1]
			path = strings.TrimSuffix(path, "/")
		}
	} else {
		fullPath = deref(it.GetName())
		path = deref(it.GetName())
	}

	if fileObj := it.GetFile(); fileObj != nil {
		// is file type
		mimeType = deref(fileObj.GetMimeType())
		ype = file.ItemTypeFile
	}

	return &file.Metadata{
		ID:         deref(it.GetId()),
		Name:       deref(it.GetName()),
		FullPath:   fullPath,
		Path:       path,
		Size:       deref(it.GetSize()),
		MimeType:   mimeType,
		ETag:       deref(it.GetETag()),
		CTag:       deref(it.GetCTag()),
		ParentID:   parentID,
		CreatedAt:  it.GetCreatedDateTime(),
		ModifiedAt: it.GetLastModifiedDateTime(),
		Type:       ype,
	}
}
