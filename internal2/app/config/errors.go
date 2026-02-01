package config

import "errors"

var (
	ErrMissingConfigurationSource = errors.New("configuration source is missing")
	// ErrAlreadyRegistered is returned when AddPath is called with a name
	// that has already been registered.
	ErrAlreadyRegistered = errors.New("configuration already registered")
	// ErrNotRegistered is returned when GetConfiguration is called for a
	// name that has no associated path.
	ErrNotRegistered = errors.New("configuration not registered")

	// ErrPathMissing is returned when a configuration name has been
	// registered but the associated path is empty or whitespace.
	ErrPathMissing = errors.New("configuration path missing")
)
