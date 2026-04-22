package logger

// Service manages the creation and configuration of Loggers.
type Service interface {
	// CreateLogger initializes a new Logger associated with the specified identifier.
	CreateLogger(id string) (Logger, error)
	// SetAllLevel sets the logging level for all managed Loggers simultaneously.
	SetAllLevel(level Level)
	// Reconfigure updates the global logging configuration, including level, output, and format.
	Reconfigure(level Level, output string, format string) error
}
