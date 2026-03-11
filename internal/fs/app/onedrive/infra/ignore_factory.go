package infra

import (
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/ignore"
)

type IgnoreMatcherFactory struct {
}

func NewIgnoreMatcherFactory() *IgnoreMatcherFactory {
	return &IgnoreMatcherFactory{}
}

func (f *IgnoreMatcherFactory) CreateMatcher(ctx context.Context, r io.Reader) (domain.IgnoreMatcher, error) {
	return ignore.ParseReader(r)
}
