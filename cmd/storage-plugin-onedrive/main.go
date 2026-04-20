package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/grpc"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/onedrive"
)

func main() {
	// Initialize the OneDrive backend.
	// PlatformProvider and other config would typically be injected or read from env.
	// For now, we instantiate the stateless backend.
	backend := onedrive.NewBackend(make(map[string]string))
	server := onedrive.NewServer(backend)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "ODC_STORAGE_PLUGIN",
			MagicCookieValue: "v1",
		},
		Plugins: map[string]plugin.Plugin{
			"onedrive": &grpc.StoragePlugin{Backend: server},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
