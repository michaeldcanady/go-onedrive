package credentialservice

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/config"
)

type ConfigurationService interface {
	GetConfiguration(ctx context.Context, name string) (config.Configuration3, error)
}
