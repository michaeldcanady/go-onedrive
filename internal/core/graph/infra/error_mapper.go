package infra

import (
	"errors"

	commonerrors "github.com/michaeldcanady/go-onedrive/internal/common/errors"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MapGraphError maps Microsoft Graph API errors to domain errors.
func MapGraphError(err error) error {
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
				return &commonerrors.DomainError{Kind: commonerrors.ErrNotFound, Err: err}

			case "accessDenied":
				return &commonerrors.DomainError{Kind: commonerrors.ErrForbidden, Err: err}

			case "unauthenticated":
				return &commonerrors.DomainError{Kind: commonerrors.ErrUnauthorized, Err: err}

			case "conflict":
				return &commonerrors.DomainError{Kind: commonerrors.ErrConflict, Err: err}

			case "preconditionFailed":
				return &commonerrors.DomainError{Kind: commonerrors.ErrPrecondition, Err: err}
			}
		}
	}

	// Try to extract HTTP status code (Kiota adapter)
	var respErr interface{ StatusCode() int }
	if errors.As(err, &respErr) {
		switch respErr.StatusCode() {
		case 401:
			return &commonerrors.DomainError{Kind: commonerrors.ErrUnauthorized, Err: err}
		case 403:
			return &commonerrors.DomainError{Kind: commonerrors.ErrForbidden, Err: err}
		case 404:
			return &commonerrors.DomainError{Kind: commonerrors.ErrNotFound, Err: err}
		case 409:
			return &commonerrors.DomainError{Kind: commonerrors.ErrConflict, Err: err}
		case 412:
			return &commonerrors.DomainError{Kind: commonerrors.ErrPrecondition, Err: err}
		case 429, 500, 502, 503, 504:
			return &commonerrors.DomainError{Kind: commonerrors.ErrTransient, Err: err}
		}
	}

	// Fallback
	return &commonerrors.DomainError{Kind: commonerrors.ErrInternal, Err: err}
}

func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}
