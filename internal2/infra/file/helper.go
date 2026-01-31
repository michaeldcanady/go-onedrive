package file

import (
	"errors"
	"path"
	"strings"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// normalizePath ensures paths like "Documents", "/Documents", "Documents/" all become "/Documents"
func normalizePath(p string) string {
	if p == "" || p == "/" || p == "." {
		return ""
	}
	return path.Clean("/" + p)
}

func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}

func toDomainItem(driveID string, it models.DriveItemable) *DriveItem {

	var mimeType string
	if file := it.GetFile(); file != nil {
		mimeType = deref(file.GetMimeType())
	}

	return &DriveItem{
		DriveID:          driveID,
		ID:               deref(it.GetId()),
		Name:             deref(it.GetName()),
		Path:             deref(it.GetParentReference().GetPath()),
		PathWithoutDrive: strings.Split(deref(it.GetParentReference().GetPath()), ":")[1],
		IsFolder:         it.GetFolder() != nil,
		Size:             deref(it.GetSize()),
		ETag:             deref(it.GetETag()),
		MimeType:         mimeType,
		Modified:         deref(it.GetLastModifiedDateTime()),
	}
}

// mapGraphError converts MS Graph OData error to domain error
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
