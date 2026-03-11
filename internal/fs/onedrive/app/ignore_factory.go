package app

import (
	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/fs/onedrive/infra"
)

func NewIgnoreMatcherFactory() domain.IgnoreMatcherFactory {
	return infra.NewIgnoreMatcherFactory()
}
