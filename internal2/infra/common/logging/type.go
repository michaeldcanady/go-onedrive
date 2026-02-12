package logging

type Type int64

const (
	TypeZap Type = iota

	DefaultLoggerType = TypeZap
)
