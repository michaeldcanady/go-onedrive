package validation

import (
	"fmt"
	"slices"
	"strings"
)

// InList creates a policy that checks if a value (or values) are within an allowed list.
func InList[T comparable](value T, allowed []T, fieldName string) error {
	if !slices.Contains(allowed, value) {
		allowedStrings := make([]string, len(allowed))
		for i, v := range allowed {
			allowedStrings[i] = fmt.Sprintf("%v", v)
		}

		return fmt.Errorf("invalid %s '%v'; please use one of the following valid options: %s",
			fieldName, value, strings.Join(allowedStrings, ", "))
	}
	return nil
}

// InListFunc returns a PolicyFunc for checking if a value is in a list.
func InListFunc[T any, V comparable](allowed []V, fieldName string, getter func(T) V) PolicyFunc[T] {
	return func(candidate T) error {
		val := getter(candidate)
		return InList(val, allowed, fieldName)
	}
}

// Required creates a policy that ensures a string field is not empty.
func Required[T any](getter func(T) string, fieldName string) PolicyFunc[T] {
	return func(candidate T) error {
		if getter(candidate) == "" {
			return fmt.Errorf("%s is required", fieldName)
		}
		return nil
	}
}
