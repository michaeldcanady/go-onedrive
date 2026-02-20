package file

import (
	"errors"
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
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

// mapGraphError converts a raw Microsoft Graph or Kiota error into a structured
// DomainError with a specific Kind classification.
//
// The function attempts multiple extraction strategies:
//
//  1. ODataError (Graph error payload)
//     - Maps well‑known Graph error codes such as:
//     "itemNotFound", "accessDenied", "conflict", "preconditionFailed", etc.
//
//  2. Kiota transport errors exposing StatusCode()
//     - Maps HTTP status codes (401, 403, 404, 409, 412, 429, 5xx).
//
//  3. Fallback
//     - Any unrecognized error becomes ErrInternal.
//
// The returned error is always a *DomainError, except when err is nil.
func mapGraphError(err error) error {
	if err == nil {
		return nil
	}

	// Try to extract OData error
	var odataErr odataerrors.ODataErrorable
	if errors.As(err, &odataErr) {
		if odataErr.GetErrorEscaped() != nil && odataErr.GetErrorEscaped().GetCode() != nil {
			code := deref(odataErr.GetErrorEscaped().GetCode())

			switch code {
			case "itemNotFound", "ErrorItemNotFound":
				return &DomainError{Kind: ErrNotFound, Err: err}

			case "accessDenied":
				return &DomainError{Kind: ErrForbidden, Err: err}

			case "unauthenticated":
				return &DomainError{Kind: ErrUnauthorized, Err: err}

			case "conflict":
				return &DomainError{Kind: ErrConflict, Err: err}

			case "preconditionFailed":
				return &DomainError{Kind: ErrPrecondition, Err: err}
			}
		}
	}

	// Try to extract HTTP status code (Kiota adapter)
	var respErr interface{ StatusCode() int }
	if errors.As(err, &respErr) {
		switch respErr.StatusCode() {
		case 401:
			return &DomainError{Kind: ErrUnauthorized, Err: err}
		case 403:
			return &DomainError{Kind: ErrForbidden, Err: err}
		case 404:
			return &DomainError{Kind: ErrNotFound, Err: err}
		case 409:
			return &DomainError{Kind: ErrConflict, Err: err}
		case 412:
			return &DomainError{Kind: ErrPrecondition, Err: err}
		case 429, 500, 502, 503, 504:
			return &DomainError{Kind: ErrTransient, Err: err}
		}
	}

	// Fallback
	return &DomainError{Kind: ErrInternal, Err: err}
}

// mapGraphError2 is a lightweight variant of mapGraphError that returns only
// domain error *kinds* rather than wrapping the original error.
//
// This is useful in contexts where callers only care about classification
// (e.g., retry logic, cache invalidation) and do not need the underlying error.
//
// Behavior mirrors mapGraphError, except:
//
//   - Known Graph/OData/HTTP errors return sentinel values such as ErrNotFound.
//   - Unknown errors return a *DomainError with Kind ErrInternal.
func mapGraphError2(err error) error {
	if err == nil {
		return nil
	}

	var odataErr odataerrors.ODataErrorable
	if errors.As(err, &odataErr) {
		if odataErr.GetErrorEscaped() != nil && odataErr.GetErrorEscaped().GetCode() != nil {
			code := deref(odataErr.GetErrorEscaped().GetCode())

			switch code {
			case "itemNotFound", "ErrorItemNotFound":
				return ErrNotFound

			case "accessDenied":
				return ErrForbidden

			case "unauthenticated":
				return ErrUnauthorized

			case "conflict":
				return ErrConflict

			case "preconditionFailed":
				return ErrPrecondition
			}
		}
	}

	// Try to extract HTTP status code (Kiota adapter)
	var respErr interface{ StatusCode() int }
	if errors.As(err, &respErr) {
		switch respErr.StatusCode() {
		case 401:
			return ErrUnauthorized
		case 403:
			return ErrForbidden
		case 404:
			return ErrNotFound
		case 409:
			return ErrConflict
		case 412:
			return ErrPrecondition
		case 429, 500, 502, 503, 504:
			return ErrTransient
		}
	}

	// Fallback
	return &DomainError{Kind: ErrInternal, Err: err}
}

// mapItemToMetadata converts a Microsoft Graph DriveItem into the internal
// file.Metadata domain model.
//
// The function extracts:
//
//   - Parent ID and normalized parent path
//   - File vs folder type (via presence of File facet)
//   - MIME type (files only)
//   - Size, ETag, CTag
//   - Created/modified timestamps
//
// Path normalization:
//
//   - Graph parent paths are of the form: "driveID:/path/to/folder"
//   - The function strips the "{driveID}:" prefix
//   - Trailing slashes are removed
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
		mimeType = *fileObj.GetMimeType()
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
