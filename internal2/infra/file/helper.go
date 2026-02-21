package file

import (
	"path"
	"strings"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// normalizePath ensures that user‑supplied paths are converted into a canonical
// OneDrive‑compatible form.
//
// The function normalizes several common path variants:
//
//   - "Documents"
//   - "/Documents"
//   - "Documents/"
//   - "./Documents"
//
// All normalize to the canonical form:
//
//	"/Documents"
//
// Special cases:
//
//   - "", "/", "." → return "" (representing the drive root)
//
// The returned value is always either "" (root) or a leading‑slash absolute path
// suitable for use with Microsoft Graph item path lookups.
func normalizePath(p string) string {
	if p == "" || p == "/" || p == "." {
		return ""
	}
	return path.Clean("/" + p)
}

// deref safely dereferences a pointer of any type T.
//
// If ptr is non‑nil, its value is returned. If ptr is nil, the zero value of T
// is returned instead. This helper eliminates repetitive nil checks when
// working with Kiota‑generated pointer‑heavy Graph models.
func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}

// toDomainItem converts a Microsoft Graph DriveItem into the internal DriveItem
// domain model.
//
// The function extracts common metadata fields, normalizes parent paths, and
// safely dereferences all pointer‑based Graph SDK fields using deref().
//
// Behavior:
//
//   - MIME type is populated only for file items.
//   - PathWithoutDrive strips the "{drive-id}:" prefix from the Graph path.
//   - IsFolder is true when the DriveItem contains a Folder facet.
//   - Missing or nil Graph fields are treated as zero values.
//
// The returned DriveItem is always non‑nil.
func toDomainItem(driveID string, it models.DriveItemable) *DriveItem {

	var mimeType string
	if file := it.GetFile(); file != nil {
		mimeType = deref(file.GetMimeType())
	}

	path := deref(it.GetParentReference().GetPath())
	pathWithoutDrive := ""
	if path != "" {
		pathWithoutDrive = strings.Split(deref(it.GetParentReference().GetPath()), ":")[1]
	}

	return &DriveItem{
		DriveID:          driveID,
		ID:               deref(it.GetId()),
		Name:             deref(it.GetName()),
		Path:             path,
		PathWithoutDrive: pathWithoutDrive,
		IsFolder:         it.GetFolder() != nil,
		Size:             deref(it.GetSize()),
		ETag:             deref(it.GetETag()),
		MimeType:         mimeType,
		Modified:         deref(it.GetLastModifiedDateTime()),
	}
}
