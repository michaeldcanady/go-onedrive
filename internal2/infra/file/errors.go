package file

import (
	"errors"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

var (
	// ErrNotFound indicates the item (file or folder) was not found.
	ErrNotFound = errors.New("not found")
	// ErrNotFolder indicates the specified item is not a folder.
	ErrNotFolder = errors.New("not a folder")
	// ErrUnauthorized indicates the user is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates the user is authenticated but does not have
	// permission for the operation.
	ErrForbidden = errors.New("forbidden")
	// ErrConflict indicates the operation failed because of a resource conflict,
	// such as a file with the same name.
	ErrConflict = errors.New("conflict")
	// ErrInternal indicates an unexpected internal error.
	ErrInternal = errors.New("internal error")
	// ErrPrecondition indicates a precondition (like ETag) check failed.
	ErrPrecondition = errors.New("precondition error")
	// ErrTransient indicates a temporary error that can be retried.
	ErrTransient = errors.New("transient")
)

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
