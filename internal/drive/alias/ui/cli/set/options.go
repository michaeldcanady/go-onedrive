package set

import "io"

// Options defines the configuration for the drive alias set operation.
type Options struct {
	// Alias is the user-friendly name to assign to the drive.
	Alias string
	// DriveID is the unique identifier of the drive to which the alias will be assigned.
	DriveID string
	// Stdout is the output writer for any command output.
	Stdout io.Writer
	// Stderr is the output writer for any command error output.
	Stderr io.Writer
}

func (o *Options) Validate() error {
	return nil
}
