package fs

import (
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

// ContainsIllegalChars checks if the path contains any characters that are not allowed in OneDrive paths.
// Returns true and an [IllegalCharacterError] if illegal characters are found. Otherwise, returns false and nil error.
func ContainsIllegalChars(path string) (bool, error) {
	illegalChars := []string{"#", "?", "*", "[", "]", "\\"}
	for _, char := range illegalChars {
		if strings.Contains(path, char) {
			return true, NewIllegalCharacterError(path, char)
		}
	}

	return false, nil
}

// ValidatePathSyntax checks for common path issues like trailing slashes or illegal characters.
func ValidatePathSyntax(p string) error {
	// Disallow trailing slashes unless it's the root path "/"
	if strings.HasSuffix(p, "/") && p != "/" {
		return NewTrailingSlashError(p)
	}

	if contains, err := ContainsIllegalChars(p); contains {
		return err
	}
	return nil
}
