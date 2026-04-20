package grpc

import (
	"context"
	"io"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/grpc/proto"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BackendMediator implements the fs.Backend and fs.AdvancedBackend interfaces.
// It acts as a client wrapper for the storage gRPC plugin.
type BackendMediator struct {
	client proto.BackendServiceClient
	plugin *plugin.Client
}

// NewBackendMediator spawns a plugin process and returns a mediator.
func NewBackendMediator(pluginPath string) (*BackendMediator, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "ODC_STORAGE_PLUGIN",
			MagicCookieValue: "v1",
		},
		Plugins: map[string]plugin.Plugin{
			"storage": &StoragePlugin{},
		},
		Cmd:              exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense("storage")
	if err != nil {
		return nil, err
	}

	return &BackendMediator{
		client: raw.(proto.BackendServiceClient),
		plugin: client,
	}, nil
}

func (b *BackendMediator) Close() {
	b.plugin.Kill()
}

func (b *BackendMediator) Name() string {
	return "grpc-plugin"
}

func (b *BackendMediator) Stat(ctx context.Context, token, driveID, path string) (fs.Item, error) {
	resp, err := b.client.Stat(ctx, &proto.StatRequest{AccessToken: token, DriveId: driveID, Path: path})
	if err != nil {
		return fs.Item{}, fromGrpcError(err)
	}
	return FromProtoItem(resp.Item), nil
}

func (b *BackendMediator) List(ctx context.Context, token, driveID, path string) ([]fs.Item, error) {
	resp, err := b.client.List(ctx, &proto.ListRequest{AccessToken: token, DriveId: driveID, Path: path})
	if err != nil {
		return nil, fromGrpcError(err)
	}

	items := make([]fs.Item, len(resp.Items))
	for i, p := range resp.Items {
		items[i] = FromProtoItem(p)
	}
	return items, nil
}

func (b *BackendMediator) Open(ctx context.Context, token, driveID, path string) (io.ReadCloser, error) {
	stream, err := b.client.Open(ctx, &proto.OpenRequest{AccessToken: token, DriveId: driveID, Path: path})
	if err != nil {
		return nil, fromGrpcError(err)
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				pw.CloseWithError(fromGrpcError(err))
				return
			}
			if _, err := pw.Write(resp.Data); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()
	return pr, nil
}

func (b *BackendMediator) Create(ctx context.Context, token, driveID, path string, r io.Reader) (fs.Item, error) {
	stream, err := b.client.Create(ctx)
	if err != nil {
		return fs.Item{}, fromGrpcError(err)
	}

	if err := stream.Send(&proto.CreateRequest{Request: &proto.CreateRequest_Path{Path: path}}); err != nil {
		return fs.Item{}, fromGrpcError(err)
	}
	if err := stream.Send(&proto.CreateRequest{Request: &proto.CreateRequest_DriveId{DriveId: driveID}}); err != nil {
		return fs.Item{}, fromGrpcError(err)
	}
	if err := stream.Send(&proto.CreateRequest{Request: &proto.CreateRequest_AccessToken{AccessToken: token}}); err != nil {
		return fs.Item{}, fromGrpcError(err)
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			if err := stream.Send(&proto.CreateRequest{Request: &proto.CreateRequest_Data{Data: buf[:n]}}); err != nil {
				return fs.Item{}, fromGrpcError(err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fs.Item{}, fromGrpcError(err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fs.Item{}, fromGrpcError(err)
	}
	return FromProtoItem(resp.Item), nil
}

func (b *BackendMediator) Mkdir(ctx context.Context, token, driveID, path string) error {
	_, err := b.client.Mkdir(ctx, &proto.MkdirRequest{AccessToken: token, DriveId: driveID, Path: path})
	return fromGrpcError(err)
}

func (b *BackendMediator) Remove(ctx context.Context, token, driveID, path string) error {
	_, err := b.client.Remove(ctx, &proto.RemoveRequest{AccessToken: token, DriveId: driveID, Path: path})
	return fromGrpcError(err)
}

func (b *BackendMediator) Capabilities() fs.Capabilities {
	resp, err := b.client.Capabilities(context.Background(), &proto.CapabilitiesRequest{})
	if err != nil {
		return fs.Capabilities{}
	}
	return FromProtoCapabilities(resp.Capabilities)
}

func (b *BackendMediator) Move(ctx context.Context, token, driveID, src, dst string) error {
	_, err := b.client.Move(ctx, &proto.MoveRequest{AccessToken: token, DriveId: driveID, Src: src, Dst: dst})
	return fromGrpcError(err)
}

func (b *BackendMediator) Copy(ctx context.Context, token, driveID, src, dst string) error {
	_, err := b.client.Copy(ctx, &proto.CopyRequest{AccessToken: token, DriveId: driveID, Src: src, Dst: dst})
	return fromGrpcError(err)
}

func fromGrpcError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return fs.ErrNotFound
	case codes.PermissionDenied:
		return fs.ErrForbidden
	case codes.AlreadyExists:
		return fs.ErrConflict
	case codes.InvalidArgument:
		return fs.ErrInvalidRequest
	default:
		return fs.ErrInternal
	}
}
