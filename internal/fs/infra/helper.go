package infra

import (
	"errors"
	"path"
	"strings"

	commonerrors "github.com/michaeldcanady/go-onedrive/internal/common/errors"
	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
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

// mapGraphError converts a raw Microsoft Graph or Kiota error into a structured
// domain error. If wrap is true, it returns a DomainError that wraps the original
// error; otherwise, it returns only the domain error kind.
func mapGraphError(err error, wrap bool) error {
	if err == nil {
		return nil
	}

	kind := commonerrors.ErrInternal

	// Try to extract OData error
	var odataErr odataerrors.ODataErrorable
	if errors.As(err, &odataErr) {
		if odataErr.GetErrorEscaped() != nil && odataErr.GetErrorEscaped().GetCode() != nil {
			code := deref(odataErr.GetErrorEscaped().GetCode())

			switch code {
			case "itemNotFound", "ErrorItemNotFound":
				kind = commonerrors.ErrNotFound
			case "accessDenied":
				kind = commonerrors.ErrForbidden
			case "unauthenticated":
				kind = commonerrors.ErrUnauthorized
			case "conflict":
				kind = commonerrors.ErrConflict
			case "preconditionFailed", "notAllowed":
				kind = commonerrors.ErrPrecondition
			case "invalidRequest":
				kind = commonerrors.ErrInvalidRequest
			}
		}
	} else {
		// Try to extract HTTP status code (Kiota adapter)
		var respErr interface{ StatusCode() int }
		if errors.As(err, &respErr) {
			switch respErr.StatusCode() {
			case 401:
				kind = commonerrors.ErrUnauthorized
			case 403:
				kind = commonerrors.ErrForbidden
			case 404:
				kind = commonerrors.ErrNotFound
			case 409:
				kind = commonerrors.ErrConflict
			case 412:
				kind = commonerrors.ErrPrecondition
			case 429, 500, 502, 503, 504:
				kind = commonerrors.ErrTransient
			}
		}
	}

	if !wrap {
		if kind == commonerrors.ErrInternal && err != nil {
			return &commonerrors.DomainError{Kind: kind, Err: err}
		}
		return kind
	}

	return &commonerrors.DomainError{Kind: kind, Err: err}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
