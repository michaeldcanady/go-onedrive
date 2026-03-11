package util

import "io"

// BaseOptions defines common configuration for all CLI commands.
type BaseOptions struct {
	Quiet bool

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
