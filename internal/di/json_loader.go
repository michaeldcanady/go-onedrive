package di

import (
	"encoding/json"
	"os"

	"github.com/michaeldcanady/go-onedrive/internal/config"
)

type JSONLoader struct{}

func (l JSONLoader) Load(path string) (config.Configuration3, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config.Configuration3{}, err
	}

	var cfg config.Configuration3
	if err := json.Unmarshal(data, &cfg); err != nil {
		return config.Configuration3{}, err
	}

	return cfg, nil
}
