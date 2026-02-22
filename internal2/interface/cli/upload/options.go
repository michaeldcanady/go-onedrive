// Package upload provides the command-line interface for uploading files to OneDrive.
package upload

import (
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// Options defines the configuration for the upload command.
// It encapsulates the source local path, the destination OneDrive path,
// and flags that control the upload behavior.
type Options struct {
	// Source is the path to the local file to upload.
	Source string

	// Destination is the path in OneDrive where the file should be uploaded.
	// If it ends with a slash, the source file's name will be appended.
	Destination string

	// Overwrite determines whether to replace an existing file at the destination.
	Overwrite bool
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if required paths are missing.
func (o *Options) Validate() error {
	if o.Source == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "source path is required")
	}
	if o.Destination == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "destination path is required")
	}
	return nil
}
