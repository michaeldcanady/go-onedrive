package domain

type LoggerProvider interface {
	Logger(opts LoggerOptions) (Logger, error)
}
