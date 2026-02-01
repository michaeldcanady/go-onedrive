package auth

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type ConfigProvider interface {
	GetConfiguration(ctx context.Context, profile string) (config.Configuration3, error)
}
