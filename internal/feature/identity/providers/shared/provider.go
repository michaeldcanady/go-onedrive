package shared

import (
	"context"

	abstractions "github.com/microsoft/kiota-abstractions-go"
)

// PlatformProvider defines the interface for a cloud service platform (e.g., Microsoft, Google).
// It serves as the primary gateway for obtaining authenticated API clients and request adapters.
type PlatformProvider interface {
	// Name returns the unique identifier for the cloud platform (e.g., "microsoft").
	Name() string
	// Adapter returns a pre-configured, authenticated Kiota request adapter for API calls.
	Adapter(ctx context.Context) (abstractions.RequestAdapter, error)
}
