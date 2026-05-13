// storage-plugin-googledrive provides a storage backend for Google Drive.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

type GoogleDriveStoragePlugin struct {
	storage_proto.UnimplementedStorageServiceServer
}

func (p *GoogleDriveStoragePlugin) ListDrives(ctx context.Context, req *storage_proto.ListDrivesRequest) (*storage_proto.ListDrivesResponse, error) {
	srv, err := p.getService(ctx, req.Options)
	if err != nil {
		return nil, err
	}
	res, err := srv.About.Get().Fields("user").Do()
	if err != nil {
		return nil, err
	}
	return &storage_proto.ListDrivesResponse{Drives: []*storage_proto.Drive{{Id: "root", Name: res.User.DisplayName + "'s Drive", Type: "personal"}}}, nil
}

func (p *GoogleDriveStoragePlugin) GetMetadata(ctx context.Context, req *storage_proto.MetadataRequest) (*storage_proto.MetadataResponse, error) {
	return &storage_proto.MetadataResponse{
		Name:               "googledrive",
		Type:               "storage",
		SupportedProviders: []string{"google"},
	}, nil
}

func (p *GoogleDriveStoragePlugin) List(ctx context.Context, req *storage_proto.ListRequest) (*storage_proto.ListResponse, error) {
	srv, err := p.getService(ctx, req.Options)
	if err != nil {
		return nil, err
	}
	id, err := p.resolvePath(srv, req.Path)
	if err != nil {
		return nil, err
	}
	res, err := srv.Files.List().Q(fmt.Sprintf("'%s' in parents and trashed = false", id)).Fields("files(id, name, mimeType, size, modifiedTime)").Do()
	if err != nil {
		return nil, err
	}
	nodes := make([]*storage_proto.Node, len(res.Files))
	for i, f := range res.Files {
		nodes[i] = p.toProtoNode(f, filepath.Join(req.Path, f.Name))
	}
	return &storage_proto.ListResponse{Nodes: nodes}, nil
}

func (p *GoogleDriveStoragePlugin) Stat(ctx context.Context, req *storage_proto.StatRequest) (*storage_proto.StatResponse, error) {
	srv, err := p.getService(ctx, req.Options)
	if err != nil {
		return nil, err
	}
	id, err := p.resolvePath(srv, req.Path)
	if err != nil {
		return nil, err
	}
	f, err := srv.Files.Get(id).Fields("id, name, mimeType, size, modifiedTime").Do()
	if err != nil {
		return nil, err
	}
	return &storage_proto.StatResponse{Node: p.toProtoNode(f, req.Path)}, nil
}

func (p *GoogleDriveStoragePlugin) Mkdir(ctx context.Context, req *storage_proto.MkdirRequest) (*storage_proto.MkdirResponse, error) {
	srv, err := p.getService(ctx, req.Options)
	if err != nil {
		return nil, err
	}
	parent, err := p.resolvePath(srv, filepath.Dir(req.Path))
	if err != nil {
		return nil, err
	}
	f, err := srv.Files.Create(&drive.File{Name: filepath.Base(req.Path), MimeType: "application/vnd.google-apps.folder", Parents: []string{parent}}).Fields("id, name, mimeType, size, modifiedTime").Do()
	if err != nil {
		return nil, err
	}
	return &storage_proto.MkdirResponse{Node: p.toProtoNode(f, req.Path)}, nil
}

func (p *GoogleDriveStoragePlugin) Read(req *storage_proto.ReadRequest, stream storage_proto.StorageService_ReadServer) error {
	srv, err := p.getService(stream.Context(), req.Options)
	if err != nil {
		return err
	}
	id, err := p.resolvePath(srv, req.Path)
	if err != nil {
		return err
	}
	res, err := srv.Files.Get(id).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf := make([]byte, 32*1024)
	for {
		n, err := res.Body.Read(buf)
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

func (p *GoogleDriveStoragePlugin) Write(stream storage_proto.StorageService_WriteServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	srv, err := p.getService(stream.Context(), req.Options)
	if err != nil {
		return err
	}
	parent, err := p.resolvePath(srv, filepath.Dir(req.Path))
	if err != nil {
		return err
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		pw.Write(req.Chunk)
		for {
			m, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}
			pw.Write(m.Chunk)
		}
	}()

	name := filepath.Base(req.Path)
	res, _ := srv.Files.List().Q(fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", name, parent)).Do()
	var f *drive.File
	if len(res.Files) > 0 {
		f, err = srv.Files.Update(res.Files[0].Id, nil).Media(pr).Fields("id, name, mimeType, size, modifiedTime").Do()
	} else {
		f, err = srv.Files.Create(&drive.File{Name: name, Parents: []string{parent}}).Media(pr).Fields("id, name, mimeType, size, modifiedTime").Do()
	}
	if err != nil {
		return err
	}
	return stream.SendAndClose(&storage_proto.WriteResponse{Node: p.toProtoNode(f, req.Path)})
}

func (p *GoogleDriveStoragePlugin) Delete(ctx context.Context, req *storage_proto.DeleteRequest) (*storage_proto.DeleteResponse, error) {
	srv, err := p.getService(ctx, req.Options)
	if err != nil {
		return nil, err
	}
	id, err := p.resolvePath(srv, req.Path)
	if err != nil {
		return nil, err
	}
	if err := srv.Files.Delete(id).Do(); err != nil {
		return nil, err
	}
	return &storage_proto.DeleteResponse{Success: true}, nil
}

func (p *GoogleDriveStoragePlugin) getService(ctx context.Context, opts map[string]string) (*drive.Service, error) {
	t := opts["token"]
	if t == "" {
		return nil, fmt.Errorf("missing token")
	}
	return drive.NewService(ctx, option.WithHTTPClient(&http.Client{Transport: &plugins.TokenTransport{Token: t}}))
}

func (p *GoogleDriveStoragePlugin) resolvePath(srv *drive.Service, path string) (string, error) {
	id := "root"
	for _, part := range strings.Split(strings.Trim(path, "/"), "/") {
		if part == "" {
			continue
		}
		res, err := srv.Files.List().Q(fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", part, id)).Fields("files(id)").Do()
		if err != nil || len(res.Files) == 0 {
			return "", fmt.Errorf("not found: %s", part)
		}
		id = res.Files[0].Id
	}
	return id, nil
}

func (p *GoogleDriveStoragePlugin) toProtoNode(f *drive.File, path string) *storage_proto.Node {
	t := storage_proto.NodeType_FILE
	if f.MimeType == "application/vnd.google-apps.folder" {
		t = storage_proto.NodeType_DIRECTORY
	}
	mod, _ := time.Parse(time.RFC3339, f.ModifiedTime)
	return &storage_proto.Node{Id: f.Id, Name: f.Name, Path: path, Type: t, Size: f.Size, ModifiedAt: mod.Unix()}
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins:         map[string]plugin.Plugin{"storage": &plugins.StorageGRPCPlugin{Impl: &GoogleDriveStoragePlugin{}}},
		GRPCServer:      plugins.CustomGRPCServer,
	})
}
