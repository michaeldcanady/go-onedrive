package logger

type Options struct {
	LogLevel          string
	LogsHome          string
	OutputDestination OutputDestination
	Type              Type
}

func NewOptions() *Options {
	return &Options{
		LogLevel:          DefaultLoggerLevel,
		OutputDestination: DefaultLoggerOutputDestination,
		Type:              DefaultLoggerType,
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

func WithType(ype Type) Option {
	return func(opts *Options) error {
		opts.Type = ype
		return nil
	}
}

func WithOutputDestinationFile(logsHome string) Option {
	return func(o *Options) error {
		o.OutputDestination = OutputDestinationFile
		o.LogsHome = logsHome
		return nil
	}
}

func WithOutputDestinationStandardOut() Option {
	return func(o *Options) error {
		o.OutputDestination = OutputDestinationStandardOut
		return nil
	}
}

func WithOutputDestinationStandardError() Option {
	return func(o *Options) error {
		o.OutputDestination = OutputDestinationStandardError
		return nil
	}
}
