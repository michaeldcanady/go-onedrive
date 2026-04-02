package registry

import (
	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
)

// Service defines the interface for managing a collection of identity providers.
type Service interface {
	// Register associates an authenticator with its provider name.
	Register(provider string, auth shared.Authenticator)
	// Get retrieves the authenticator for a given provider name.
	Get(provider string) (shared.Authenticator, error)
}
