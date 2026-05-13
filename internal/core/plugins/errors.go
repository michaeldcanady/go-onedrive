package plugins

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/michaeldcanady/go-onedrive/internal/core/errors"
)

// FromGRPC translates a gRPC error into a standard [errors] domain error.
// If the error is not a gRPC status error, it is returned as-is.
func FromGRPC(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.ErrNotFound
	case codes.AlreadyExists:
		return errors.ErrAlreadyExists
	case codes.PermissionDenied, codes.Unauthenticated:
		return errors.ErrPermissionDenied
	case codes.InvalidArgument:
		return errors.ErrInvalidPath
	case codes.Internal:
		return errors.ErrInternal
	case codes.Unavailable:
		return errors.ErrUnavailable
	case codes.FailedPrecondition:
		// Often used for directory not empty or similar
		if strings.Contains(st.Message(), "not empty") {
			return errors.ErrNotEmpty
		}
		return err
	default:
		return err
	}
}
