package root

import "errors"

var (
	ErrMissingConfigSection = errors.New("missing config section")
	ErrUnmarshalConfig      = errors.New("unable to unmarshal config")
	ErrCreateCredential     = errors.New("unable to create credential")
	ErrMissingAuthConfig    = errors.New("missing 'auth' config section")
	ErrMissingLoggingConfig = errors.New("missing 'logging' config section")
	ErrUnmarshalAuthConfig  = errors.New("unable to unmarshal auth config")
	ErrUnmarshalLogConfig   = errors.New("unable to unmarshal logging config")
)
