package edit

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// Name returns the filename without its extension.
func Name(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// isAuthRequired checks if an error indicates that authentication is needed.
func isAuthRequired(err error) bool {
	var authErr *azidentity.AuthenticationRequiredError
	return errors.As(err, &authErr)
}
