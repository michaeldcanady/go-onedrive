package fs

import "fmt"

// IllegalCharacterError represents an error where a path contains characters that are not allowed in OneDrive paths.
type IllegalCharacterError struct {
	Path      string
	Character string
}

// NewIllegalCharacterError creates a new IllegalCharacterError for the given path and character.
func NewIllegalCharacterError(path, character string) *IllegalCharacterError {
	return &IllegalCharacterError{
		Path:      path,
		Character: character,
	}
}

func (e *IllegalCharacterError) Error() string {
	return fmt.Sprintf("path %s contains illegal character %s", e.Path, e.Character)
}
