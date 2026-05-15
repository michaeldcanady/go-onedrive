// storage-plugin-onedrive provides a storage backend for Microsoft OneDrive.
package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-plugin"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphdrives "github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

type OneDriveStoragePlugin struct {
	storage_proto.UnimplementedStorageServiceServer
}

func (p *OneDriveStoragePlugin) List(ctx context.Context, req *storage_proto.ListRequest) (*storage_proto.ListResponse, error) {
	c, err := p.getClient(req.Options)
	if err != nil {
		return nil, err
	}
	res, err := c.Drives().ByDriveId(p.getDriveID(req.Options)).Items().ByDriveItemId(p.resolvePath(req.Path)).Children().Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	nodes := make([]*storage_proto.Node, 0)
	for _, item := range res.GetValue() {
		name := ""
		if item.GetName() != nil {
			name = *item.GetName()
		}
		nodes = append(nodes, p.toProtoNode(item, filepath.Join(req.Path, name)))
	}
	return &storage_proto.ListResponse{Nodes: nodes}, nil
}

func (p *OneDriveStoragePlugin) Stat(ctx context.Context, req *storage_proto.StatRequest) (*storage_proto.StatResponse, error) {
	c, err := p.getClient(req.Options)
	if err != nil {
		return nil, err
	}
	item, err := c.Drives().ByDriveId(p.getDriveID(req.Options)).Items().ByDriveItemId(p.resolvePath(req.Path)).Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &storage_proto.StatResponse{Node: p.toProtoNode(item, req.Path)}, nil
}

func (p *OneDriveStoragePlugin) Mkdir(ctx context.Context, req *storage_proto.MkdirRequest) (*storage_proto.MkdirResponse, error) {
	c, err := p.getClient(req.Options)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(req.Path)
	f := models.NewDriveItem()
	f.SetName(&name)
	f.SetFolder(models.NewFolder())
	item, err := c.Drives().ByDriveId(p.getDriveID(req.Options)).Items().ByDriveItemId(p.resolvePath(filepath.Dir(req.Path))).Children().Post(ctx, f, nil)
	if err != nil {
		return nil, err
	}
	return &storage_proto.MkdirResponse{Node: p.toProtoNode(item, req.Path)}, nil
}

func (p *OneDriveStoragePlugin) Read(req *storage_proto.ReadRequest, stream storage_proto.StorageService_ReadServer) error {
	c, err := p.getClient(req.Options)
	if err != nil {
		return err
	}
	b, err := c.Drives().ByDriveId(p.getDriveID(req.Options)).Items().ByDriveItemId(p.resolvePath(req.Path)).Content().Get(stream.Context(), nil)
	if err != nil {
		return err
	}
	return stream.Send(&storage_proto.ReadResponse{Chunk: b})
}

func (p *OneDriveStoragePlugin) Write(stream storage_proto.StorageService_WriteServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	c, err := p.getClient(req.Options)
	if err != nil {
		return err
	}

	var data []byte
	data = append(data, req.Chunk...)
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		data = append(data, m.Chunk...)
	}

	var cfg *msgraphdrives.ItemItemsItemContentRequestBuilderPutRequestConfiguration
	if etag := req.Options["if_match"]; etag != "" {
		cfg = &msgraphdrives.ItemItemsItemContentRequestBuilderPutRequestConfiguration{Headers: abstractions.NewRequestHeaders()}
		cfg.Headers.Add("If-Match", etag)
	}

	item, err := c.Drives().ByDriveId(p.getDriveID(req.Options)).Items().ByDriveItemId(p.resolvePath(req.Path)).Content().Put(stream.Context(), data, cfg)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&storage_proto.WriteResponse{Node: p.toProtoNode(item, req.Path)})
}

func (p *OneDriveStoragePlugin) Delete(ctx context.Context, req *storage_proto.DeleteRequest) (*storage_proto.DeleteResponse, error) {
	c, err := p.getClient(req.Options)
	if err != nil {
		return nil, err
	}
	if err := c.Drives().ByDriveId(p.getDriveID(req.Options)).Items().ByDriveItemId(p.resolvePath(req.Path)).Delete(ctx, nil); err != nil {
		return nil, err
	}
	return &storage_proto.DeleteResponse{Success: true}, nil
}

func (p *OneDriveStoragePlugin) ListDrives(ctx context.Context, req *storage_proto.ListDrivesRequest) (*storage_proto.ListDrivesResponse, error) {
	c, err := p.getClient(req.Options)
	if err != nil {
		return nil, err
	}
	res, err := c.Me().Drives().Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	drives := make([]*storage_proto.Drive, 0)
	for _, d := range res.GetValue() {
		drives = append(drives, &storage_proto.Drive{Id: *d.GetId(), Name: *d.GetName(), Type: *d.GetDriveType()})
	}
	return &storage_proto.ListDrivesResponse{Drives: drives}, nil
}

func (p *OneDriveStoragePlugin) GetMetadata(ctx context.Context, req *storage_proto.MetadataRequest) (*storage_proto.MetadataResponse, error) {
	return &storage_proto.MetadataResponse{
		Name:               "onedrive",
		Type:               "storage",
		SupportedProviders: []string{"azure"},
	}, nil
}

func (p *OneDriveStoragePlugin) getClient(opts map[string]string) (*msgraph.GraphServiceClient, error) {
	t := opts["token"]
	if t == "" {
		return nil, fmt.Errorf("missing token")
	}
	tp := &plugins.TokenTransport{Token: t}
	adapter, err := msgraph.NewGraphRequestAdapter(&authProvider{tp})
	if err != nil {
		return nil, err
	}
	return msgraph.NewGraphServiceClient(adapter), nil
}

type authProvider struct{ tp *plugins.TokenTransport }

func (a *authProvider) AuthenticateRequest(ctx context.Context, req *abstractions.RequestInformation, _ map[string]any) error {
	if req.Headers == nil {
		req.Headers = abstractions.NewRequestHeaders()
	}
	req.Headers.Add("Authorization", "Bearer "+a.tp.Token)
	return nil
}

func (p *OneDriveStoragePlugin) getDriveID(opts map[string]string) string {
	if id := opts["drive_id"]; id != "" {
		return id
	}
	return "root"
}

func (p *OneDriveStoragePlugin) resolvePath(path string) string {
	if path == "" || path == "/" {
		return "root"
	}
	return "root:/" + strings.TrimPrefix(path, "/") + ":"
}

func (p *OneDriveStoragePlugin) toProtoNode(item models.DriveItemable, path string) *storage_proto.Node {
	node := &storage_proto.Node{Path: path, Type: storage_proto.NodeType_FILE}
	if item.GetFolder() != nil {
		node.Type = storage_proto.NodeType_DIRECTORY
	}
	if s := item.GetSize(); s != nil {
		node.Size = *s
	}
	if t := item.GetLastModifiedDateTime(); t != nil {
		node.ModifiedAt = t.Unix()
	}
	if e := item.GetETag(); e != nil {
		node.Etag = *e
	}
	if c := item.GetCTag(); c != nil {
		node.Ctag = *c
	}
	if i := item.GetId(); i != nil {
		node.Id = *i
	}
	if n := item.GetName(); n != nil {
		node.Name = *n
	}
	return node
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins:         map[string]plugin.Plugin{"storage": &plugins.StorageGRPCPlugin{Impl: &OneDriveStoragePlugin{}}},
		GRPCServer:      plugins.CustomGRPCServer,
	})
}
