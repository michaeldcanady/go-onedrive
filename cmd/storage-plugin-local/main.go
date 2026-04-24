package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/michaeldcanady/go-onedrive/internal/features/storage/backend/grpc"
	"github.com/michaeldcanady/go-onedrive/internal/features/storage/backend/local"
)

func main() {
	// Initialize the local backend (stateless).
	// Root should ideally be passed via env var or flag in a real plugin.
	backend := local.NewBackend("/tmp/odc-storage")
	server := local.NewServer(backend)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "ODC_STORAGE_PLUGIN",
			MagicCookieValue: "v1",
		},
		Plugins: map[string]plugin.Plugin{
			"local": &grpc.StoragePlugin{Backend: server},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
