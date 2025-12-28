package cmd

import "errors"

var (
	ErrMissingConfigSection = errors.New("missing config section")
	ErrUnmarshalConfig      = errors.New("unable to unmarshal config")
	ErrCreateCredential     = errors.New("unable to create credential")
)
