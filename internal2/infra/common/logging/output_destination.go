package logging

import "strings"

type OutputDestination int64

const (
	OutputDestinationUnknown OutputDestination = iota - 1
	OutputDestinationStandardOut
	OutputDestinationFile
	OutputDestinationStandardError

	DefaultLoggerOutputDestination = OutputDestinationFile
	DefaultLoggerLevel             = "info"
)

var outputDestinationNames = map[OutputDestination]string{
	OutputDestinationUnknown:       "unknown",
	OutputDestinationStandardOut:   "stdout",
	OutputDestinationFile:          "file",
	OutputDestinationStandardError: "stderr",
}

var outputDestinationValues = map[string]OutputDestination{
	"unknown": OutputDestinationUnknown,
	"stdout":  OutputDestinationStandardOut,
	"file":    OutputDestinationFile,
	"stderr":  OutputDestinationStandardError,
}

func ParseOutputDestination(s string) OutputDestination {
	if dest, ok := outputDestinationValues[strings.ToLower(s)]; ok {
		return dest
	}
	return OutputDestinationUnknown
}

func (o OutputDestination) String() string {
	if name, ok := outputDestinationNames[o]; ok {
		return name
	}
	return OutputDestinationUnknown.String()
}
