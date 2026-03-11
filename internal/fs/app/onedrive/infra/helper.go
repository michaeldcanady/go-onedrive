package infra

import (
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

// normalizePath ensures that user‑supplied paths are converted into a canonical
// OneDrive‑compatible form.
func normalizePath(p string) string {
	if p == "" || p == "/" || p == "." {
		return ""
	}
	return path.Clean("/" + p)
}

// deref safely dereferences a pointer of any type T.
func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}

// mapItemToMetadata converts a Microsoft Graph DriveItem into the internal
// domain.Metadata domain model.
func mapItemToMetadata(it models.DriveItemable) *domain.Metadata {
	if it == nil {
		return nil
	}

	var (
		parentID string
		mimeType string
		fullPath string
		path     string
		itemType domain.ItemType = domain.ItemTypeFolder
	)

	if parent := it.GetParentReference(); parent != nil {
		parentID = deref(parent.GetId())
		fullPath = deref(parent.GetPath())
		if fullPath != "" {
			// Graph paths are "drives/{drive-id}/root:/path/to/item"
			parts := strings.Split(fullPath, ":")
			if len(parts) > 1 {
				path = parts[1]
				path = strings.TrimSuffix(path, "/")
			}
		}
	} else {
		fullPath = deref(it.GetName())
		path = deref(it.GetName())
	}

	if fileObj := it.GetFile(); fileObj != nil {
		mimeType = deref(fileObj.GetMimeType())
		itemType = domain.ItemTypeFile
	}

	return &domain.Metadata{
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
		Type:       itemType,
	}
}

// expandPathTemplate expands a URI template for a drive item, choosing between
// root and relative templates based on whether the path is empty.
func expandPathTemplate(rootTemplate, relativeTemplate, driveID, path string) string {
	normalized := normalizePath(path)
	urlTemplate := rootTemplate
	subs := make(stduritemplate.Substitutions)
	subs["baseurl"] = baseURL
	subs["drive_id"] = driveID

	if normalized != "" {
		urlTemplate = relativeTemplate
		subs["path"] = normalized
	}

	uri, _ := stduritemplate.Expand(urlTemplate, subs)
	return uri
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
