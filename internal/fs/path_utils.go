package fs

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
)

// ContainsIllegalChars checks if the path contains any characters that are not allowed in OneDrive paths.
// Returns true and an [IllegalCharacterError] if illegal characters are found. Otherwise, returns false and nil error.
func ContainsIllegalChars(path string) (bool, error) {
	illegalChars := []string{"#", "?", "*", "[", "]", "\\"}
	for _, char := range illegalChars {
		if strings.Contains(path, char) {
			return true, errors.NewIllegalCharacterError(path, char, nil)
		}
	}

	return false, nil
}

// ValidatePathSyntax checks for common path issues like trailing slashes or illegal characters.
func ValidatePathSyntax(p string) error {
	// Disallow trailing slashes unless it's the root path "/"
	if strings.HasSuffix(p, "/") && p != "/" {
		return errors.NewTrailingSlashError(p, nil)
	}

	if contains, err := ContainsIllegalChars(p); contains {
		return err
	}
	return nil
}
