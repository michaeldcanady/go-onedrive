package infra

import (
	"os"

	domainconfig "github.com/michaeldcanady/go-onedrive/internal/config/domain"
	"github.com/stretchr/testify/assert/yaml"
)

type YAMLLoader struct{}

func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

func (l *YAMLLoader) Load(path string) (domainconfig.Configuration, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domainconfig.Configuration{}, err
	}

	var cfg domainconfig.Configuration
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return domainconfig.Configuration{}, err
	}

	return cfg, nil
}
