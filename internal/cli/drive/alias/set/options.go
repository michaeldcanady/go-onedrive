package set

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
)

// Options defines the configuration for the drive alias set command.
// It encapsulates the alias name, target drive ID, and standard I/O streams.
type Options struct {
	// Alias is the friendly name to assign to the drive.
	Alias string

	// DriveID is the unique identifier of the OneDrive drive.
	DriveID string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if either Alias or DriveID is missing.
func (o *Options) Validate() error {
	if o.Alias == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "alias is required")
	}
	if o.DriveID == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "drive ID is required")
	}
	return nil
}
