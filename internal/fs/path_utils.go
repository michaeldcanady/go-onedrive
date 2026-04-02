package fs

import (
	"fmt"
	"strings"
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

// ValidatePathSyntax checks for common path issues like trailing slashes or illegal characters.
func ValidatePathSyntax(p string) error {
	// Disallow trailing slashes unless it's the root path "/"
	if strings.HasSuffix(p, "/") && p != "/" {
		return fmt.Errorf("path '%s' has a trailing slash, which is not allowed", p)
	}

	// Disallow illegal characters (e.g., '#', '?', '*', '[', ']', '\') - Windows illegal chars
	// Note: We exclude ':' here because it's handled as a provider/alias separator by SplitProviderPath.
	// However, if the path *after* the provider still contains ':', it might be invalid depending on the provider.
	illegalChars := []string{"#", "?", "*", "[", "]", "\\"}
	for _, char := range illegalChars {
		if strings.Contains(p, char) {
			return fmt.Errorf("path '%s' contains illegal character '%s'", p, char)
		}
	}

	return nil
}
