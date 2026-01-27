package ls

import (
	"fmt"

	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
)

type OutputFormatter interface {
	Format(items []*driveservice.DriveItem) error
}

type FormatterFactory struct{}

func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{}
}

func (f *FormatterFactory) Create(format string) (OutputFormatter, error) {
	switch format {
	case "short", "":
		return &HumanShortFormatter{}, nil

	case "long":
		return &HumanLongFormatter{}, nil

	case "json":
		return &JSONFormatter{}, nil

	case "yaml", "yml":
		return &YAMLFormatter{}, nil

	default:
		return nil, fmt.Errorf("invalid format: %s", format)
	}
}

type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(items []*driveservice.DriveItem) error {
	return printYAML(items)
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(items []*driveservice.DriveItem) error {
	return printJSON(items)
}

type HumanLongFormatter struct{}

func (f *HumanLongFormatter) Format(items []*driveservice.DriveItem) error {
	printLongDomain(items)
	return nil
}

type HumanShortFormatter struct{}

func (f *HumanShortFormatter) Format(items []*driveservice.DriveItem) error {
	printShortDomain(items)
	return nil
}
