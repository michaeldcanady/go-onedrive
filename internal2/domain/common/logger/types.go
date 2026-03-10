package logger

import "strings"

type Type int

const (
	TypeUnknown Type = iota
	TypeZap
)

const (
	DefaultLoggerType = TypeZap
)

type OutputDestination int

const (
	OutputDestinationUnknown OutputDestination = iota
	OutputDestinationStandardOut
	OutputDestinationStandardError
	OutputDestinationFile
)

const (
	DefaultLoggerOutputDestination = OutputDestinationStandardOut
)

const (
	DefaultLoggerLevel = "info"
)

type LoggerOptions struct {
	Level             string
	OutputDestination OutputDestination
	FilePath          string
}

type LoggerOption func(*LoggerOptions) error

func NewLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		Level:             DefaultLoggerLevel,
		OutputDestination: DefaultLoggerOutputDestination,
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

func WithLoggerLogLevel(level string) LoggerOption {
	return func(o *LoggerOptions) error {
		o.Level = level
		return nil
	}
}

func WithLoggerOutputDestinationStandardOut() LoggerOption {
	return func(o *LoggerOptions) error {
		o.OutputDestination = OutputDestinationStandardOut
		return nil
	}
}

func WithLoggerOutputDestinationStandardError() LoggerOption {
	return func(o *LoggerOptions) error {
		o.OutputDestination = OutputDestinationStandardError
		return nil
	}
}

func WithLoggerOutputDestinationFile(path string) LoggerOption {
	return func(o *LoggerOptions) error {
		o.OutputDestination = OutputDestinationFile
		o.FilePath = path
		return nil
	}
}

func ParseOutputDestination(s string) OutputDestination {
	switch strings.ToLower(s) {
	case "stdout", "standardout":
		return OutputDestinationStandardOut
	case "stderr", "standarderror":
		return OutputDestinationStandardError
	case "file":
		return OutputDestinationFile
	default:
		return OutputDestinationUnknown
	}
}
