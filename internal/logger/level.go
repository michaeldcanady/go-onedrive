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
	// LevelFatal represents fatal error messages.
	LevelFatal
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
	case "fatal":
		return LevelFatal
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
	case LevelFatal:
		return "fatal"
	case LevelUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// MarshalText implements encoding.TextMarshaler for serializing Level as a string.
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for deserializing Level from a string.
func (l *Level) UnmarshalText(text []byte) error {
	*l = ParseLevel(string(text))
	return nil
}
