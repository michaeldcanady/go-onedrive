package logging

type LoggerOptions struct {
	Level             string
	outputDestination OutputDestination
	logPath           string
}

func NewLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		Level:             "info",
		outputDestination: OutputDestinationStandardOut,
		logPath:           "",
	}
}

func (o *LoggerOptions) Apply(opts ...LoggerOption) error {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}
	return nil
}

type LoggerOption func(*LoggerOptions) error

func WithLogLevel(level string) LoggerOption {
	return func(o *LoggerOptions) error {
		o.Level = level
		return nil
	}
}

func WithOutputDestinationFile(logPath string) LoggerOption {
	return func(o *LoggerOptions) error {
		o.outputDestination = OutputDestinationFile
		o.logPath = logPath
		return nil
	}
}

func WithOutputDestinationStandardOut() LoggerOption {
	return func(o *LoggerOptions) error {
		o.outputDestination = OutputDestinationStandardOut
		return nil
	}
}

func WithOutputDestinationStandardError() LoggerOption {
	return func(o *LoggerOptions) error {
		o.outputDestination = OutputDestinationStandardError
		return nil
	}
}
