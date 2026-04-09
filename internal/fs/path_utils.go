package fs

import (
	"errors"
	"fmt"
	"strings"

	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
)

// SplitProviderPath splits a path into its provider prefix and the remaining path.
// For example, "local:/etc/hosts" returns "local", "/etc/hosts", true.
// "my-alias:/some/path" returns "my-alias", "/some/path", true.
// "/plain/path" returns "", "/plain/path", false.
func SplitProviderPath(path string) (string, string, bool) {
	prefix, rest, found := strings.Cut(path, ":")
	if !found {
		return "", path, false
	}
	return prefix, rest, true
}

func ContainsIllegalChars(path string) bool {
	// Disallow illegal characters (e.g., '#', '?', '*', '[', ']', '\') - Windows illegal chars
	// Note: We exclude ':' here because it's handled as a provider/alias separator by SplitProviderPath.
	// However, if the path *after* the provider still contains ':', it might be invalid depending on the provider.
	illegalChars := []string{"#", "?", "*", "[", "]", "\\"}
	for _, char := range illegalChars {
		if strings.Contains(path, char) {
			return true
		}
	}

	return false
}

// ValidatePathSyntax checks for common path issues like trailing slashes or illegal characters.
func ValidatePathSyntax(p string) error {
	// Disallow trailing slashes unless it's the root path "/"
	if strings.HasSuffix(p, "/") && p != "/" {
		return coreerrors.NewInvalidInput(
			errors.New("path has a trailing slash, which is not allowed"),
			fmt.Sprintf("invalid path '%s' due to trailing slash", p),
			"Remove the trailing slash from the path",
		)
	}

	if ContainsIllegalChars(p) {
		return coreerrors.NewInvalidInput(
			errors.New("path contains illegal characters"),
			fmt.Sprintf("invalid path '%s' due to illegal characters", p),
			"Remove the illegal characters from the path",
		)
	}
	return nil
}
