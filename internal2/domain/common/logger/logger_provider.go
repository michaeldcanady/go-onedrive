package logger

type LoggerProvider interface {
	Logger(opts LoggerOptions) (Logger, error)
}
