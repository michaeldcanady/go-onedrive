package config

import (
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

type YAMLLoader struct{}

func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

func (l *YAMLLoader) Load(path string) (Configuration3, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Configuration3{}, err
	}

	var cfg Configuration3
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Configuration3{}, err
	}

	return cfg, nil
}
