package infra

import (
	"errors"

	commonerrors "github.com/michaeldcanady/go-onedrive/internal/common/errors"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MapGraphError maps Microsoft Graph API errors to domain errors.
// If wrap is true, it always returns a DomainError that wraps the original error.
// If wrap is false, it returns the domain error variable itself, unless the error
// is an internal error, in which case it always returns a DomainError.
func MapGraphError(err error, wrap bool) error {
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

func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}
