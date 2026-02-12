package logger

import "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"

type Options struct {
	LogLevel          string
	LogsHome          string
	OutputDestination logging.OutputDestination
	Type              logging.Type
}

func NewOptions() *Options {
	return &Options{
		LogLevel:          logging.DefaultLoggerLevel,
		OutputDestination: logging.DefaultLoggerOutputDestination,
		Type:              logging.DefaultLoggerType,
	}
}

func (o *Options) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}
	return nil
}

type Option func(*Options) error

func WithLogLevel(level string) Option {
	return func(opts *Options) error {
		opts.LogLevel = level
		return nil
	}
}

func WithType(ype logging.Type) Option {
	return func(opts *Options) error {
		opts.Type = ype
		return nil
	}
}

func WithOutputDestinationFile(logsHome string) Option {
	return func(o *Options) error {
		o.OutputDestination = logging.OutputDestinationFile
		o.LogsHome = logsHome
		return nil
	}
}

func WithOutputDestinationStandardOut() Option {
	return func(o *Options) error {
		o.OutputDestination = logging.OutputDestinationStandardOut
		return nil
	}
}

func WithOutputDestinationStandardError() Option {
	return func(o *Options) error {
		o.OutputDestination = logging.OutputDestinationStandardError
		return nil
	}
}
