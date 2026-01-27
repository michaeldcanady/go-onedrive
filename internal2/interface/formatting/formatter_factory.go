package formatting

import (
	"fmt"

	domainformatting "github.com/michaeldcanady/go-onedrive/internal2/domain/formatting"
)

type FormatterFactory struct{}

func NewFormatterFactory() *FormatterFactory { return &FormatterFactory{} }

func (f *FormatterFactory) Create(format string) (domainformatting.OutputFormatter, error) {
	if format == "" {
		format = "short"
	}
	switch format {
	case "short", "":
		return &HumanShortFormatter{}, nil
	case "long":
		return &HumanLongFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "yaml", "yml":
		return &YAMLFormatter{}, nil
	}
	return nil, fmt.Errorf("invalid format: %s", format)
}
