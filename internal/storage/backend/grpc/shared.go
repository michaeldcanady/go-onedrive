package grpc

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/grpc/proto"
	"google.golang.org/grpc"
)

// StoragePlugin is the implementation of plugin.GRPCPlugin.
type StoragePlugin struct {
	plugin.Plugin
	// Backend is the server-side implementation.
	Backend proto.BackendServiceServer
}

func (p *StoragePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterBackendServiceServer(s, p.Backend)
	return nil
}

func (p *StoragePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewBackendServiceClient(c), nil
}
