package logger

import "github.com/michaeldcanady/go-onedrive/pkg/logger"

type Logger = logger.Logger
type Field = logger.Field
type FieldType = logger.FieldType
type Level = logger.Level

const (
	LevelDebug   = logger.LevelDebug
	LevelInfo    = logger.LevelInfo
	LevelWarn    = logger.LevelWarn
	LevelError   = logger.LevelError
	LevelUnknown = logger.LevelUnknown
)

type _int interface {
	int | int32 | int64
}

func Int[T _int](key string, val T) Field {
	return logger.Int(key, val)
}

var (
	ParseLevel = logger.ParseLevel
	String     = logger.String
	Time       = logger.Time
	Duration   = logger.Duration
	Error      = logger.Error
	Bool       = logger.Bool
)
