package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

type MetadataCache interface {
	Get(ctx context.Context, id string) (*file.Metadata, bool)
	Put(ctx context.Context, m *file.Metadata) error
	Invalidate(ctx context.Context, id string) error
}
