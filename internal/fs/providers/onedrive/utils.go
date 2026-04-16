package onedrive

import (
	"context"
	"errors"
	"path"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

func mapError(err error, path string) error {
	if err == nil {
		return nil
	}

	kind := fs.ErrInternal
	var apiErr *abstractions.ApiError
	if errors.As(err, &apiErr) {
		switch apiErr.ResponseStatusCode {
		case 401:
			// mapping unauthorized to forbidden for simplicity in pkg/fs for now,
			// or we could add ErrUnauthorized to pkg/fs/errors.go
			kind = fs.ErrForbidden
		case 403:
			kind = fs.ErrForbidden
		case 404:
			kind = fs.ErrNotFound
		case 409:
			kind = fs.ErrConflict
		}
	}

	return &fs.Error{
		Kind: kind,
		Err:  err,
		Path: path,
	}
}

func resolveDriveID(ctx context.Context, uri *fs.URI, resolver fs.DriveResolver) string {
	if uri.DriveID != "" {
		return uri.DriveID
	}

	if resolver != nil {
		driveID, err := resolver.GetActiveDriveID(ctx)
		if err == nil && driveID != "" {
			return driveID
		}
	}

	return "me"
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
