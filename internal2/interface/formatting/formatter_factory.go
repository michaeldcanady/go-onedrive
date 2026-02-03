package formatting

import (
	"fmt"

	domainformatting "github.com/michaeldcanady/go-onedrive/internal2/domain/formatting"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type FormatterFactory struct{}

func NewFormatterFactory() *FormatterFactory { return &FormatterFactory{} }

func (f *FormatterFactory) Create(format string) (domainformatting.OutputFormatter[domainfs.Item], error) {
	if format == "" {
		format = "short"
	}
	switch format {
	case "short", "":
		return &HumanShortFormatter{term: Terminal{}}, nil
	case "long":
		return &HumanLongFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "yaml", "yml":
		return &YAMLFormatter{}, nil
	case "tree":
		return NewTreeFormatter(), nil
	}
	return nil, fmt.Errorf("invalid format: %s", format)
}
