package driveservice2

import (
	"errors"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

type DriveType string

const (
	DriveTypePersonal   DriveType = "personal"
	DriveTypeBusiness   DriveType = "business"
	DriveTypeSharePoint DriveType = "sharepoint"
)

type Drive struct {
	ID       string
	Name     string
	Type     DriveType
	Owner    string
	ReadOnly bool
}

func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}

func toDomainDrive(g models.Driveable) *Drive {
	if g == nil {
		return nil
	}

	return &Drive{
		ID:   deref(g.GetId()),
		Name: deref(g.GetName()),
		Type: DriveType(deref(g.GetDriveType())),
	}
}

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
