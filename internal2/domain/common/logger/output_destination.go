package logger

import "strings"

type OutputDestination int

const (
	OutputDestinationUnknown OutputDestination = iota
	OutputDestinationStandardOut
	OutputDestinationStandardError
	OutputDestinationFile
)

func ParseOutputDestination(s string) OutputDestination {
	switch strings.ToLower(s) {
	case "stdout", "standardout":
		return OutputDestinationStandardOut
	case "stderr", "standarderror":
		return OutputDestinationStandardError
	case "file":
		return OutputDestinationFile
	default:
		return OutputDestinationUnknown
	}
}
