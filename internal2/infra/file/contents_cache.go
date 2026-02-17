package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

type ContentsCache interface {
	Get(ctx context.Context, id string) (*file.Contents, bool)
	Put(ctx context.Context, id string, m *file.Contents) error
	Invalidate(ctx context.Context, id string) error
}
