package plugins

import (
	"context"
	"net/http"

	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

// Metadata represents the self-described capabilities and identity of a plugin.
type Metadata struct {
	Name               string
	Type               string
	SupportedProviders []string
	PluginPath         string
}

// Manager coordinates the lifecycle and communication of external plugins.
// It handles the discovery, launching, and gRPC connection management for plugin processes.
type Manager interface {
	// GetStoragePlugin returns a gRPC client for the specified storage plugin.
	GetStoragePlugin(name string) (storage_proto.StorageServiceClient, error)

	// GetIdentityPlugin returns a gRPC client for the specified identity plugin.
	GetIdentityPlugin(name string) (identity_proto.IdentityPluginClient, error)

	// ListPlugins returns a list of all discovered plugins and their metadata.
	ListPlugins(ctx context.Context) ([]*Metadata, error)

	// Shutdown terminates all active plugin processes and cleans up associated resources.
	Shutdown(ctx context.Context) error
}

// TokenTransport is an [http.RoundTripper] that injects a bearer token into the Authorization header of every request.
type TokenTransport struct {
	Token string
	Base  http.RoundTripper
}

// RoundTrip executes a single HTTP transaction, adding the bearer token to the request headers.
// If [TokenTransport.Base] is nil, [http.DefaultTransport] is used.
func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.Token)

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(req)
}
