package cliprofileservicego

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/config"
)

type ConfigurationService interface {
	AddPath(name, path string) error
	GetConfiguration(ctx context.Context, name string) (config.Configuration3, error)
}
