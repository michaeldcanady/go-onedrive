package auth

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type ConfigurationService interface {
	GetConfiguration(ctx context.Context, name string) (config.Configuration3, error)
}
