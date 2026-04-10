package microsoft

import "errors"

var (
	// ErrNotAuthenticated is returned when an operation is attempted without valid authentication credentials.
	ErrNotAuthenticated = errors.New("no authentication credential provided; please run 'login' first")
)
