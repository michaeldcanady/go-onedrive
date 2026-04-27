package logger

import (
	"encoding"
	"fmt"
	"strconv"
	"strings"
)

// Level defines the severity of a log message.
type Level int8

var _ encoding.TextMarshaler = Level(0)
var _ encoding.TextUnmarshaler = (*Level)(nil)

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

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(text []byte) error {
	s := string(text)
	*l = ParseLevel(s)
	if *l != LevelUnknown {
		return nil
	}

	// Try parsing as integer for backward compatibility
	if i, err := strconv.Atoi(s); err == nil {
		if i >= -128 && i <= 127 {
			lvl := Level(i)
			switch lvl {
			case LevelDebug, LevelInfo, LevelWarn, LevelError:
				*l = lvl
				return nil
			}
		}
	}

	return fmt.Errorf("invalid log level: %s", s)
}

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
