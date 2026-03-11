package domain

import "io"

// Launcher defines the interface for launching an external editordomain.
type Launcher interface {
	Launch(path string) error
	WithIO(stdin io.Reader, stdout, stderr io.Writer) Launcher
}

// Service defines the interface for editor-related operations.
type Service interface {
	LaunchTempFile(prefix, suffix string, reader io.Reader) ([]byte, string, error)
	WithIO(stdin io.Reader, stdout, stderr io.Writer) Service
}
