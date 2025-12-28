package logging

// Logger is an interface for logging messages with various severity levels.
type Logger interface {
	// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Info(msg string, kv ...Field)
	// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Warn(msg string, kv ...Field)
	// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Error(msg string, kv ...Field)
	// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Debug(msg string, kv ...Field)
}
