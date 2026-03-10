package fs

import (
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/ignore"
)

type IgnoreMatcherFactory struct {
}

func NewIgnoreMatcherFactory() *IgnoreMatcherFactory {
	return &IgnoreMatcherFactory{}
}

func (f *IgnoreMatcherFactory) CreateMatcher(ctx context.Context, r io.Reader) (fs.IgnoreMatcher, error) {
	return ignore.ParseReader(r)
}
