package util

import "strings"

// ParsePath separates a path into its provider prefix and the actual path within that provider.
// If no prefix is found, it defaults to "onedrive".
func ParsePath(path string) (string, string) {
	prefix, rest, found := strings.Cut(path, ":")
	if !found {
		return "onedrive", path
	}

	rest = strings.TrimPrefix(rest, "//")

	return strings.ToLower(prefix), rest
}
