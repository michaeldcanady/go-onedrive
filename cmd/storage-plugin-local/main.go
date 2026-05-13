// storage-plugin-local provides a storage backend for the local host filesystem.
package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-plugin"

	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

type LocalStoragePlugin struct {
	storage_proto.UnimplementedStorageServiceServer
}

func (p *LocalStoragePlugin) getPath(opts map[string]string, path string) string {
	root := opts["root_path"]
	if root == "" {
		root = "."
	}
	return filepath.Join(root, path)
}

func (p *LocalStoragePlugin) List(ctx context.Context, req *storage_proto.ListRequest) (*storage_proto.ListResponse, error) {
	es, err := os.ReadDir(p.getPath(req.Options, req.Path))
	if err != nil {
		return nil, err
	}
	nodes := make([]*storage_proto.Node, 0, len(es))
	for _, e := range es {
		if info, err := e.Info(); err == nil {
			nodes = append(nodes, p.toProtoNode(info, filepath.Join(req.Path, e.Name())))
		}
	}
	return &storage_proto.ListResponse{Nodes: nodes}, nil
}

func (p *LocalStoragePlugin) Stat(ctx context.Context, req *storage_proto.StatRequest) (*storage_proto.StatResponse, error) {
	info, err := os.Stat(p.getPath(req.Options, req.Path))
	if err != nil {
		return nil, err
	}
	return &storage_proto.StatResponse{Node: p.toProtoNode(info, req.Path)}, nil
}

func (p *LocalStoragePlugin) Mkdir(ctx context.Context, req *storage_proto.MkdirRequest) (*storage_proto.MkdirResponse, error) {
	full := p.getPath(req.Options, req.Path)
	if err := os.MkdirAll(full, 0755); err != nil {
		return nil, err
	}
	info, _ := os.Stat(full)
	return &storage_proto.MkdirResponse{Node: p.toProtoNode(info, req.Path)}, nil
}

func (p *LocalStoragePlugin) Read(req *storage_proto.ReadRequest, stream storage_proto.StorageService_ReadServer) error {
	f, err := os.Open(p.getPath(req.Options, req.Path))
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, 32*1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&storage_proto.ReadResponse{Chunk: buf[:n]}); err != nil {
			return err
		}
	}
	return nil
}

func (p *LocalStoragePlugin) Write(stream storage_proto.StorageService_WriteServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	full := p.getPath(req.Options, req.Path)
	os.MkdirAll(filepath.Dir(full), 0755)
	f, err := os.Create(full)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(req.Chunk)
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f.Write(m.Chunk)
	}
	info, _ := os.Stat(full)
	return stream.SendAndClose(&storage_proto.WriteResponse{Node: p.toProtoNode(info, req.Path)})
}

func (p *LocalStoragePlugin) Delete(ctx context.Context, req *storage_proto.DeleteRequest) (*storage_proto.DeleteResponse, error) {
	if err := os.RemoveAll(p.getPath(req.Options, req.Path)); err != nil {
		return nil, err
	}
	return &storage_proto.DeleteResponse{Success: true}, nil
}

func (p *LocalStoragePlugin) Move(ctx context.Context, req *storage_proto.MoveRequest) (*storage_proto.MoveResponse, error) {
	src, dst := p.getPath(req.Options, req.Source), p.getPath(req.Options, req.Destination)
	os.MkdirAll(filepath.Dir(dst), 0755)
	if err := os.Rename(src, dst); err != nil {
		return nil, err
	}
	info, _ := os.Stat(dst)
	return &storage_proto.MoveResponse{Node: p.toProtoNode(info, req.Destination)}, nil
}

func (p *LocalStoragePlugin) ListDrives(ctx context.Context, req *storage_proto.ListDrivesRequest) (*storage_proto.ListDrivesResponse, error) {
	return &storage_proto.ListDrivesResponse{
		Drives: []*storage_proto.Drive{
			{Id: "/", Name: "Local Filesystem", Type: "local"},
		},
	}, nil
}

func (p *LocalStoragePlugin) GetMetadata(ctx context.Context, req *storage_proto.MetadataRequest) (*storage_proto.MetadataResponse, error) {
	return &storage_proto.MetadataResponse{
		Name:               "local",
		Type:               "storage",
		SupportedProviders: []string{},
	}, nil
}

func (p *LocalStoragePlugin) toProtoNode(info os.FileInfo, path string) *storage_proto.Node {
	t := storage_proto.NodeType_FILE
	if info.IsDir() {
		t = storage_proto.NodeType_DIRECTORY
	}
	return &storage_proto.Node{Name: info.Name(), Path: path, Type: t, Size: info.Size(), ModifiedAt: info.ModTime().Unix()}
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins:         map[string]plugin.Plugin{"storage": &plugins.StorageGRPCPlugin{Impl: &LocalStoragePlugin{}}},
		GRPCServer:      plugins.CustomGRPCServer,
	})
}
