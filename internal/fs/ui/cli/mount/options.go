package mount

import (
	"fmt"
	"io"
	"strings"
)

// Options defines the configuration for the mount operation.
type Options struct {
	// Path is the virtual path where the backend will be mounted.
	Path string
	// Type is the type of backend (e.g., "local", "onedrive").
	Type string
	// IdentityID is the specific account to use for this mount.
	IdentityID string
	// RawOptions is the list of key=value pairs from the CLI.
	RawOptions []string
	// MountOptions contains backend-specific settings.
	MountOptions map[string]string
	// Stdout is the destination for the operation's output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Path == "" {
		return fmt.Errorf("mount path is required")
	}
	if !strings.HasPrefix(o.Path, "/") {
		return fmt.Errorf("mount path must be absolute (start with /)")
	}
	if o.Type == "" {
		return fmt.Errorf("backend type is required")
	}

	o.MountOptions = make(map[string]string)
	for _, opt := range o.RawOptions {
		k, v, found := strings.Cut(opt, "=")
		if !found {
			return fmt.Errorf("invalid option format %s, expected key=value", opt)
		}
		o.MountOptions[k] = v
	}

	return nil
}
