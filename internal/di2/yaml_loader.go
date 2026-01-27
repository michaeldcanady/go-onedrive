package di2

import (
	"os"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/stretchr/testify/assert/yaml"
)

type YAMLLoader struct{}

func (l YAMLLoader) Load(path string) (config.Configuration3, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config.Configuration3{}, err
	}

	var cfg config.Configuration3
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return config.Configuration3{}, err
	}

	return cfg, nil
}
