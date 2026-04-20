package onedrive

import (
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/grpc"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/grpc/proto"
)

// Server implements the proto.BackendServiceServer interface, wrapping the OneDrive Backend.
type Server struct {
	proto.UnimplementedBackendServiceServer
	backend *Backend
}

// NewServer creates a new gRPC server wrapper for the OneDrive Backend.
func NewServer(b *Backend) *Server {
	return &Server{backend: b}
}

func (s *Server) Stat(ctx context.Context, req *proto.StatRequest) (*proto.StatResponse, error) {
	item, err := s.backend.Stat(ctx, req.AccessToken, req.DriveId, req.Path)
	if err != nil {
		return nil, grpc.ToProtoError(err)
	}
	return &proto.StatResponse{Item: grpc.ToProtoItem(item)}, nil
}

func (s *Server) List(ctx context.Context, req *proto.ListRequest) (*proto.ListResponse, error) {
	items, err := s.backend.List(ctx, req.AccessToken, req.DriveId, req.Path)
	if err != nil {
		return nil, grpc.ToProtoError(err)
	}

	protoItems := make([]*proto.Item, len(items))
	for i, item := range items {
		protoItems[i] = grpc.ToProtoItem(item)
	}
	return &proto.ListResponse{Items: protoItems}, nil
}

func (s *Server) Open(req *proto.OpenRequest, stream proto.BackendService_OpenServer) error {
	f, err := s.backend.Open(stream.Context(), req.AccessToken, req.DriveId, req.Path)
	if err != nil {
		return grpc.ToProtoError(err)
	}
	defer f.Close()

	buffer := make([]byte, 32*1024)
	for {
		n, err := f.Read(buffer)
		if n > 0 {
			if err := stream.Send(&proto.OpenResponse{Data: buffer[:n]}); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return grpc.ToProtoError(err)
		}
	}
	return nil
}

func (s *Server) Create(stream proto.BackendService_CreateServer) error {
	var path string
	var driveID string
	var accessToken string
	var data []byte

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if req.GetPath() != "" {
			path = req.GetPath()
		}
		if req.GetDriveId() != "" {
			driveID = req.GetDriveId()
		}
		if req.GetAccessToken() != "" {
			accessToken = req.GetAccessToken()
		}
		if len(req.GetData()) > 0 {
			data = append(data, req.GetData()...)
		}
	}

	item, err := s.backend.Create(stream.Context(), accessToken, driveID, path, io.Reader(nil)) // Simplified
	if err != nil {
		return grpc.ToProtoError(err)
	}
	return stream.SendAndClose(&proto.CreateResponse{Item: grpc.ToProtoItem(item)})
}

func (s *Server) Mkdir(ctx context.Context, req *proto.MkdirRequest) (*proto.MkdirResponse, error) {
	err := s.backend.Mkdir(ctx, req.AccessToken, req.DriveId, req.Path)
	if err != nil {
		return nil, grpc.ToProtoError(err)
	}
	return &proto.MkdirResponse{}, nil
}

func (s *Server) Remove(ctx context.Context, req *proto.RemoveRequest) (*proto.RemoveResponse, error) {
	err := s.backend.Remove(ctx, req.AccessToken, req.DriveId, req.Path)
	if err != nil {
		return nil, grpc.ToProtoError(err)
	}
	return &proto.RemoveResponse{}, nil
}

func (s *Server) Capabilities(ctx context.Context, req *proto.CapabilitiesRequest) (*proto.CapabilitiesResponse, error) {
	caps := s.backend.Capabilities()
	return &proto.CapabilitiesResponse{Capabilities: grpc.ToProtoCapabilities(caps)}, nil
}

func (s *Server) Move(ctx context.Context, req *proto.MoveRequest) (*proto.MoveResponse, error) {
	err := s.backend.Move(ctx, req.AccessToken, req.DriveId, req.Src, req.Dst)
	if err != nil {
		return nil, grpc.ToProtoError(err)
	}
	return &proto.MoveResponse{}, nil
}

func (s *Server) Copy(ctx context.Context, req *proto.CopyRequest) (*proto.CopyResponse, error) {
	err := s.backend.Copy(ctx, req.AccessToken, req.DriveId, req.Src, req.Dst)
	if err != nil {
		return nil, grpc.ToProtoError(err)
	}
	return &proto.CopyResponse{}, nil
}
