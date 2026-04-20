package grpc

import (
	"time"

	proto "github.com/michaeldcanady/go-onedrive/internal/storage/backend/grpc/proto"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToProtoItem converts an fs.Item to the proto.Item message.
func ToProtoItem(item fs.Item) *proto.Item {
	return &proto.Item{
		Id:         item.ID,
		Name:       item.Name,
		Path:       item.Path,
		Type:       ToProtoItemType(item.Type),
		Size:       item.Size,
		ModifiedAt: item.ModifiedAt.Unix(),
		Etag:       item.ETag,
	}
}

// ToProtoItemType converts an fs.ItemType to the proto.ItemType enum.
func ToProtoItemType(t fs.ItemType) proto.ItemType {
	switch t {
	case fs.TypeFile:
		return proto.ItemType_ITEM_TYPE_FILE
	case fs.TypeFolder:
		return proto.ItemType_ITEM_TYPE_FOLDER
	default:
		return proto.ItemType_ITEM_TYPE_UNKNOWN
	}
}

// FromProtoItem converts a proto.Item message to an fs.Item struct.
func FromProtoItem(p *proto.Item) fs.Item {
	return fs.Item{
		ID:         p.Id,
		Name:       p.Name,
		Path:       p.Path,
		Type:       FromProtoItemType(p.Type),
		Size:       p.Size,
		ModifiedAt: time.Unix(p.ModifiedAt, 0),
		ETag:       p.Etag,
	}
}

// FromProtoItemType converts a proto.ItemType enum to an fs.ItemType.
func FromProtoItemType(t proto.ItemType) fs.ItemType {
	switch t {
	case proto.ItemType_ITEM_TYPE_FILE:
		return fs.TypeFile
	case proto.ItemType_ITEM_TYPE_FOLDER:
		return fs.TypeFolder
	default:
		return fs.TypeUnknown
	}
}

// ToProtoCapabilities converts an fs.Capabilities to the proto.Capabilities message.
func ToProtoCapabilities(c fs.Capabilities) *proto.Capabilities {
	return &proto.Capabilities{
		CanMove:     c.CanMove,
		CanCopy:     c.CanCopy,
		CanRecursive: c.CanRecursive,
	}
}

// FromProtoCapabilities converts a proto.Capabilities message to an fs.Capabilities struct.
func FromProtoCapabilities(p *proto.Capabilities) fs.Capabilities {
	return fs.Capabilities{
		CanMove:     p.CanMove,
		CanCopy:     p.CanCopy,
		CanRecursive: p.CanRecursive,
	}
}

// ToProtoError maps fs.Error to gRPC status codes.
func ToProtoError(err error) error {
	if err == nil {
		return nil
	}

	var fsErr *fs.Error
	// If it's a fs.Error, map the Kind
	if _, ok := err.(*fs.Error); ok {
		fsErr = err.(*fs.Error)
	}

	code := codes.Internal
	if fsErr != nil {
		switch fsErr.Kind {
		case fs.ErrNotFound:
			code = codes.NotFound
		case fs.ErrForbidden:
			code = codes.PermissionDenied
		case fs.ErrConflict:
			code = codes.AlreadyExists
		case fs.ErrInvalidRequest:
			code = codes.InvalidArgument
		}
	}

	return status.Errorf(code, err.Error())
}
