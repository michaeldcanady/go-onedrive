package file

import (
	"errors"
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
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

func mapItemToMetadata(it models.DriveItemable) *file.Metadata {
	var (
		parentID string
		mimeType string
		path     string
	)
	if parent := it.GetParentReference(); parent != nil {
		parentID = *parent.GetId()
		path = *parent.GetPath()
	}

	if file := it.GetFile(); file != nil {
		// is file type
		mimeType = *file.GetMimeType()
	}

	return &file.Metadata{
		ID:         *it.GetId(),
		Name:       *it.GetName(),
		Path:       path,
		Size:       *it.GetSize(),
		MimeType:   mimeType,
		ETag:       *it.GetETag(),
		CTag:       *it.GetCTag(),
		ParentID:   parentID,
		CreatedAt:  it.GetCreatedDateTime(),
		ModifiedAt: it.GetLastModifiedDateTime(),
	}
}
