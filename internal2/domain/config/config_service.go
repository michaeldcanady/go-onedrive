package config

import (
	"context"
)

type ConfigService interface {
	GetConfiguration(ctx context.Context, name string) (Configuration, error)
	AddPath(name, path string) error
}
