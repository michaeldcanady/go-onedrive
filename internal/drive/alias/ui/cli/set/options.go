package set

import "io"

// Options defines the configuration for the drive alias set operation.
type Options struct {
	Alias   string
	DriveID string
	Stdout  io.Writer
}
