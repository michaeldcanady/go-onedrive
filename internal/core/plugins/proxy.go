package plugins

import (
	"context"

	"google.golang.org/grpc"

	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

type storageProxy struct {
	manager *pluginManager
	name    string
}

func (p *storageProxy) client() (storage_proto.StorageServiceClient, error) {
	return p.manager.getRawStoragePlugin(p.name)
}

func (p *storageProxy) List(ctx context.Context, in *storage_proto.ListRequest, opts ...grpc.CallOption) (*storage_proto.ListResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.List(ctx, in, opts...)
}

func (p *storageProxy) Stat(ctx context.Context, in *storage_proto.StatRequest, opts ...grpc.CallOption) (*storage_proto.StatResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Stat(ctx, in, opts...)
}

func (p *storageProxy) Mkdir(ctx context.Context, in *storage_proto.MkdirRequest, opts ...grpc.CallOption) (*storage_proto.MkdirResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Mkdir(ctx, in, opts...)
}

func (p *storageProxy) Read(ctx context.Context, in *storage_proto.ReadRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[storage_proto.ReadResponse], error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Read(ctx, in, opts...)
}

func (p *storageProxy) Write(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[storage_proto.WriteRequest, storage_proto.WriteResponse], error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Write(ctx, opts...)
}

func (p *storageProxy) Delete(ctx context.Context, in *storage_proto.DeleteRequest, opts ...grpc.CallOption) (*storage_proto.DeleteResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Delete(ctx, in, opts...)
}

func (p *storageProxy) Move(ctx context.Context, in *storage_proto.MoveRequest, opts ...grpc.CallOption) (*storage_proto.MoveResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Move(ctx, in, opts...)
}

func (p *storageProxy) ListDrives(ctx context.Context, in *storage_proto.ListDrivesRequest, opts ...grpc.CallOption) (*storage_proto.ListDrivesResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.ListDrives(ctx, in, opts...)
}

func (p *storageProxy) GetDrive(ctx context.Context, in *storage_proto.GetDriveRequest, opts ...grpc.CallOption) (*storage_proto.GetDriveResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.GetDrive(ctx, in, opts...)
}

func (p *storageProxy) GetMetadata(ctx context.Context, in *storage_proto.MetadataRequest, opts ...grpc.CallOption) (*storage_proto.MetadataResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.GetMetadata(ctx, in, opts...)
}

type identityProxy struct {
	manager *pluginManager
	name    string
}

func (p *identityProxy) client() (identity_proto.IdentityPluginClient, error) {
	return p.manager.getRawIdentityPlugin(p.name)
}

func (p *identityProxy) Login(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[identity_proto.LoginRequest, identity_proto.LoginResponse], error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Login(ctx, opts...)
}

func (p *identityProxy) Refresh(ctx context.Context, in *identity_proto.RefreshRequest, opts ...grpc.CallOption) (*identity_proto.RefreshResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Refresh(ctx, in, opts...)
}

func (p *identityProxy) ListIdentities(ctx context.Context, in *identity_proto.ListIdentitiesRequest, opts ...grpc.CallOption) (*identity_proto.ListIdentitiesResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.ListIdentities(ctx, in, opts...)
}

func (p *identityProxy) Logout(ctx context.Context, in *identity_proto.LogoutRequest, opts ...grpc.CallOption) (*identity_proto.LogoutResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.Logout(ctx, in, opts...)
}

func (p *identityProxy) GetMetadata(ctx context.Context, in *identity_proto.MetadataRequest, opts ...grpc.CallOption) (*identity_proto.MetadataResponse, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}
	return c.GetMetadata(ctx, in, opts...)
}
