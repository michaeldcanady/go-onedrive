package app

import (
	"github.com/michaeldcanady/go-onedrive/internal/fs/onedrive/infra"
	"github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
)

func NewIgnoreMatcherFactory() domain.IgnoreMatcherFactory {
	return infra.NewIgnoreMatcherFactory()
}
