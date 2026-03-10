package logger

type Type int

const (
	TypeUnknown Type = iota
	TypeZap
)

const (
	DefaultLoggerType = TypeZap
)

const (
	DefaultLoggerOutputDestination = OutputDestinationStandardOut
)

const (
	DefaultLoggerLevel = "info"
)
