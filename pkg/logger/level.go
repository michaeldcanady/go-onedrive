package logger

import "strings"

// Level defines the severity of a log message.
type Level int8

const (
	// LevelUnknown represents an unknown logging level.
	LevelUnknown Level = iota - 1
	// LevelDebug represents detailed information messages.
	LevelDebug
	// LevelInfo represents general informational messages.
	LevelInfo
	// LevelWarn represents warning messages.
	LevelWarn
	// LevelError represents error messages.
	LevelError
)

func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "unknown", "":
		return LevelUnknown
	default:
		return LevelUnknown
	}
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}
