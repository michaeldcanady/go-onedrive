package onedrive

import (
	"errors"
	"path"

	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

func mapError(err error, path string) error {
	if err == nil {
		return nil
	}

	kind := coreerrors.ErrInternal
	var apiErr *abstractions.ApiError
	if errors.As(err, &apiErr) {
		switch apiErr.ResponseStatusCode {
		case 401:
			kind = coreerrors.ErrUnauthorized
		case 403:
			kind = coreerrors.ErrForbidden
		case 404:
			kind = coreerrors.ErrNotFound
		case 409:
			kind = coreerrors.ErrConflict
		case 412:
			kind = coreerrors.ErrPrecondition
		case 429, 503, 504:
			kind = coreerrors.ErrTransient
		}
	}

	return &coreerrors.DomainError{
		Kind: kind,
		Err:  err,
		Path: path,
	}
}

func resolveDriveID(uri *shared.URI, stateSvc state.Service) string {
	if uri.DriveID != "" {
		return uri.DriveID
	}

	// Default to active drive or "me"
	driveID, err := stateSvc.Get(state.KeyDrive)
	if err != nil || driveID == "" {
		// Fallback to primary drive
		return "me"
	}

	return driveID
}

func expandURI(rootTemplate, relativeTemplate, driveID, itemPath string) string {
	normalized := normalizePath(itemPath)
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

func normalizePath(pth string) string {
	if pth == "" || pth == "/" || pth == "." {
		return ""
	}
	// OneDrive relative paths in URI templates usually start with /
	return path.Clean("/" + pth)
}
