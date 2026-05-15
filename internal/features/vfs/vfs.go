package vfs

import (
	"context"
	"io"

	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

// Node represents a file or directory entry within the Virtual File System.
// It encapsulates metadata required for both display and optimistic concurrency control.
type Node struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Type       NodeType `json:"type"`
	Size       int64    `json:"size"`
	ModifiedAt int64    `json:"modified_at"`
	ETag       string   `json:"etag"` // ETag used for optimistic concurrency in Write operations.
	CTag       string   `json:"ctag"`
}

// NodeType distinguishes between files and directories.
type NodeType int

const (
	FileType      NodeType = 0
	DirectoryType NodeType = 1
)

// VFS coordinates file operations across a unified virtual namespace.
// It is the primary interface for all filesystem-like interactions in the application.
type VFS interface {
	// List returns the immediate children of the directory at the specified path.
	List(ctx context.Context, path string) ([]*Node, error)

	// Stat retrieves metadata for the file or directory at the specified path.
	Stat(ctx context.Context, path string) (*Node, error)

	// Mkdir creates a new directory at the specified path.
	Mkdir(ctx context.Context, path string) error

	// Remove deletes the file or directory at the specified path.
	Remove(ctx context.Context, path string) error

	// Move renames or relocates a node. Cross-mount moves are handled as copy-then-delete.
	Move(ctx context.Context, src, dst string) error

	// Copy replicates a node. Cross-mount copies involve streaming data through the host.
	Copy(ctx context.Context, src, dst string) error

	// Read opens a stream for reading the file's content.
	// The caller is responsible for closing the returned [io.ReadCloser].
	Read(ctx context.Context, path string) (io.ReadCloser, error)

	// Write streams data to the specified path, creating or overwriting the file.
	Write(ctx context.Context, path string, reader io.Reader, options ...WriteOption) error
}

// WriteOption configures the behavior of a [VFS.Write] operation.
type WriteOption func(map[string]string)

// WithIfMatch enables optimistic concurrency by requiring the file's ETag to match the provided value.
func WithIfMatch(etag string) WriteOption {
	return func(opts map[string]string) {
		if etag != "" {
			opts["if_match"] = etag
		}
	}
}

func FromProtoNode(p *storage_proto.Node) *Node {
	return &Node{
		ID:         p.Id,
		Name:       p.Name,
		Path:       p.Path,
		Type:       NodeType(p.Type),
		Size:       p.Size,
		ModifiedAt: p.ModifiedAt,
		ETag:       p.Etag,
		CTag:       p.Ctag,
	}
}
