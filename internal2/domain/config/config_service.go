package config

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type ConfigService interface {
	GetConfiguration(ctx context.Context, name string) (config.Configuration3, error)
	AddPath(name, path string) error
}
